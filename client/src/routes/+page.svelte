<script lang="ts">
	import { DotsHorizontalOutline, PlusOutline, SortOutline } from 'flowbite-svelte-icons';
	import * as Card from '$lib/components/ui/card/index';
	import { flip } from 'svelte/animate';
	import { dndzone } from 'svelte-dnd-action';
	import Button from '$lib/components/ui/button/button.svelte';
	import TicketCard from '$lib/components/ticket-card/ticket-card.svelte';
	import type { Unsubscriber } from 'svelte/motion';
	import { onDestroy, onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { TicketService } from '$lib/services/ticket-service';
	import { BoardStore } from '$lib/stores/board';
	import { AlertStore } from '$lib/stores/alert';
	import { DialogStore } from '$lib/stores/dialog';
	import BoardSaveDialogContent from '$lib/components/board-save-dialog-content/board-save-dialog-content.svelte';
	import StatusCreateDialogContent from '$lib/components/status-save-dialog-content/status-save-dialog-content.svelte';
	import * as Scroll from '$lib/components/ui/scroll-area';
	import TicketSaveDialogContent from '$lib/components/ticket-save-dialog-content/ticket-save-dialog-content.svelte';
	import { cloneDeep } from 'lodash';
	import {
		Subject,
		Subscription,
		catchError,
		concatMap,
		from,
		map,
		of,
		tap,
		buffer,
		debounceTime
	} from 'rxjs';
	import StatusSaveDialogContent from '$lib/components/status-save-dialog-content/status-save-dialog-content.svelte';
	import { AuthStore } from '$lib/stores/auth';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import StatusCard from '$lib/components/status-card/status-card.svelte';

	const { ScrollArea } = Scroll;

	let unsubscribers: Unsubscriber[] = [];
	let boardState = BoardStore.defaultState();
	let tempBoardState = BoardStore.defaultState();
	const ticketSubject = new Subject<{
		id: number;
		ticket_ids: number[];
	}>();
	let tickets$: Subscription;

	const releaseBufferSubject = new Subject<void>();
	onMount(async () => {
		unsubscribers.push(AuthStore.Use());
		tickets$ = ticketSubject
			.pipe(
				map((value) => {
					setTimeout(() => {
						releaseBufferSubject.next();
					}, 10);

					return value;
				}),
				buffer(releaseBufferSubject.pipe(debounceTime(1000))),
				tap((values) => {
					console.log('buffer', values);
				}),
				concatMap((values) =>
					from(
						TicketService.bulkUpdateTicketOrderInStatuses(
							{ board_id: boardState.selected.id },
							{ statuses: values }
						)
					).pipe(
						map(({ data }) => {
							return data;
						}),
						catchError((error) => {
							AlertStore.error(error);

							return of();
						})
					)
				)
			)
			.subscribe((values) => {
				if (values.length === 0) return;
				console.log('values', values);
				BoardStore.update((state) => {
					const selected = cloneDeep(state.selected);
					for (const { id, tickets } of values) {
						const idx = selected.statuses.findIndex((status) => status.id === id);

						if (idx > -1) {
							selected.statuses[idx].tickets = tickets;
						}
					}

					state.selected = selected;

					return state;
				});
			});

		unsubscribers.push(
			BoardStore.subscribe(async (state) => {
				console.log('board state change', state);

				boardState = cloneDeep(state);
				tempBoardState = cloneDeep(state);
			})
		);

		try {
			const { data: boards } = await TicketService.getBoards();

			if (boards.length > 0) {
				const { data: selected } = await TicketService.getBoardById(boards[0].id);

				BoardStore.selectBoard(selected);
			}

			BoardStore.update((state) => {
				state.boards = boards;
				return state;
			});
		} catch ({ error, message }: any) {
			AlertStore.create({
				title: 'Error',
				message: error ? error.message : message || 'An error occurred'
			});
		}
	});

	onDestroy(() => {
		unsubscribers.forEach((unsubscriber) => unsubscriber());
		console.log('unsubscribes');
		tickets$.unsubscribe();
	});

	async function fetchBoardFullDetail(board: TicketService.Board) {
		try {
			const { data: boardFullDetail } = await TicketService.getBoardById(board.id);

			BoardStore.selectBoard(boardFullDetail);
		} catch (error: any) {
			AlertStore.error(error);
		}
	}

	const flipDurationMs = 200;

	type ColumnEvent = CustomEvent & { detail: { items: TicketService.Status[] } };
	function handleDndConsiderColumns(e: ColumnEvent) {
		if (!boardState.selected) return;
		boardState.selected.statuses = [...e.detail.items];
	}
	async function handleDndFinalizeColumns(e: ColumnEvent) {
		if (!boardState.selected) return;
		boardState.selected.statuses = [...e.detail.items];
		try {
			const { data: statuses } = await TicketService.updateStatusesSortOrder(
				{
					board_id: boardState.selected.id
				},
				{ statuses: boardState.selected.statuses }
			);

			BoardStore.update((state) => {
				state.selected.statuses = statuses as TicketService.Status[];

				return state;
			});
		} catch (error) {
			AlertStore.error(error);

			BoardStore.selectBoard(tempBoardState.selected);
		}
	}

	type CardEvent = CustomEvent & { detail: { items: TicketService.Ticket[] } };
	function handleDndConsiderCards(cid: number, e: CardEvent) {
		if (!boardState.selected) return;
		const colIdx = boardState.selected.statuses?.findIndex((c) => c.id === cid);
		boardState.selected.statuses[colIdx].tickets = e.detail.items;
		boardState.selected.statuses = [...boardState.selected.statuses];
	}
	async function handleDndFinalizeCards(cid: number, e: CardEvent) {
		if (!boardState.selected) return;
		const colIdx = boardState.selected.statuses.findIndex((c) => c.id === cid);
		boardState.selected.statuses[colIdx].tickets = e.detail.items;
		boardState.selected.statuses = [...boardState.selected.statuses];

		console.log('handleDndFinalizeCards');
		ticketSubject.next({
			id: cid,
			ticket_ids: boardState.selected.statuses[colIdx].tickets.map((t) => t.id)
		});
	}

	function onSortTicketsInStatus(status: TicketService.Status) {
		let idx = boardState.selected.statuses.findIndex((s) => s.id === status.id);

		if (idx === -1) return;

		boardState.selected.statuses[idx] = status;

		ticketSubject.next({
			id: status.id,
			ticket_ids: status.tickets.map((t) => t.id)
		});
	}
</script>

<section class="grid flex-1 grid-cols-12">
	<div class="relative col-span-2 pl-2 bg-muted">
		<Button
			variant="ghost"
			class="absolute top-0 right-0 flex items-center justify-center h-auto p-1 m-1.5 hover:bg-opacity-10 hover:bg-accent-foreground"
			on:click={() => {
				DialogStore.create({ component: BoardSaveDialogContent });
			}}
		>
			<PlusOutline class="size-4" />
		</Button>
		<ScrollArea
			class="mt-8 mr-2 w-full h-[calc(100vh-(var(--header-height)+var(--footer-height))-2rem-1px)] [&>[data-melt-scroll-area-thumb]]:bg-red-400"
			orientation="vertical"
			scrollbarYClasses="dark:[&>[data-melt-scroll-area-thumb]]:bg-primary-foreground"
		>
			<div class="flex flex-col w-full gap-2">
				{#each boardState.boards as board (board.id)}
					<Button
						variant="ghost"
						class="justify-between py-2 hover:bg-opacity-10 hover:bg-accent-foreground group/sidemenu {boardState.selected &&
						boardState.selected.id === board.id
							? 'bg-accent-foreground bg-opacity-10 text-accent-foreground'
							: ''}"
						on:click={() => fetchBoardFullDetail(board)}
					>
						{board.title}
						<Button
							variant="ghost"
							class="flex items-center justify-center invisible h-auto p-1 hover:bg-opacity-10 hover:bg-accent-foreground group-hover/sidemenu:visible"
							on:click={() => {
								DialogStore.create({ component: BoardSaveDialogContent, params: { model: board } });
							}}
						>
							<DotsHorizontalOutline class="size-4" />
						</Button>
					</Button>
				{/each}
			</div>
		</ScrollArea>
	</div>
	<div class="h-full col-span-10">
		<ScrollArea orientation="horizontal" class="has-[>div>div>div]:h-full">
			{#if boardState.selected}
				<div
					class="flex justify-start gap-4 p-4 overflow-x-auto overflow-y-hidden"
					use:dndzone={{
						items: boardState.selected.statuses,
						flipDurationMs,
						type: 'columns',
						dropTargetStyle: {}
					}}
					on:consider={handleDndConsiderColumns}
					on:finalize={handleDndFinalizeColumns}
				>
					{#each boardState.selected.statuses as status, i (status)}
						<div animate:flip={{ duration: flipDurationMs }}>
							<ScrollArea
								orientation="vertical"
								class="h-[calc(100vh-(var(--header-height)+var(--footer-height))-2rem-1px)]"
							>
								<StatusCard {status} {onSortTicketsInStatus}>
									{#if status.tickets.length === 0}
										<Button
											variant="ghost"
											class="flex items-center justify-start w-full h-auto p-1 rounded hover:bg-accent"
											on:click={() => {
												DialogStore.create({
													component: TicketSaveDialogContent,
													params: { board_id: boardState?.selected?.id, status_id: status.id }
												});
											}}
										>
											<PlusOutline class="size-4" />
											<span class="ml-2">Add ticket</span>
										</Button>
									{/if}
									<div
										class="absolute top-0 left-0 flex flex-col w-full gap-2 border border-red-600 h-"
										class:absolute={status.tickets.length === 0}
										class:h-32={status.tickets.length === 0}
										use:dndzone={{
											items: status.tickets,
											flipDurationMs,
											dropTargetStyle: {}
										}}
										on:consider={(e) => handleDndConsiderCards(status.id, e)}
										on:finalize={(e) => handleDndFinalizeCards(status.id, e)}
									>
										{#each status.tickets as ticket (ticket.id)}
											<div
												tabindex={ticket.id}
												role="button"
												animate:flip={{ duration: flipDurationMs }}
											>
												<TicketCard {ticket} board_id={status.board_id} />
											</div>
										{/each}
									</div>
								</StatusCard>
							</ScrollArea>
						</div>
					{/each}
					<div class="block">
						<button
							class="flex items-center justify-start p-2 mt-4 rounded min-w-80 hover:bg-accent"
							on:click={() => {
								DialogStore.create({
									component: StatusCreateDialogContent,
									params: { board_id: boardState?.selected?.id }
								});
							}}
						>
							<PlusOutline class="size-4" />
							<span class="ml-2">Add status</span>
						</button>
					</div>
				</div>
			{/if}
		</ScrollArea>
	</div>
</section>
