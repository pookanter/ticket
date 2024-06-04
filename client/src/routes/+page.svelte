<script lang="ts">
	import * as BoardTabs from '$lib/components/board-tabs/index';
	import { DotsHorizontalOutline, PlusOutline } from 'flowbite-svelte-icons';
	import * as Card from '$lib/components/ui/card/index';
	import { flip } from 'svelte/animate';
	import { dndzone } from 'svelte-dnd-action';
	import Button from '$lib/components/ui/button/button.svelte';
	import TicketCard from '$lib/components/ticket-card/ticket-card.svelte';
	import type { Unsubscriber } from 'svelte/motion';
	import { onDestroy, onMount } from 'svelte';
	import authStore from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { TicketService } from '$lib/services/ticket-service';
	import { BoardStore, type BoardState } from '$lib/stores/board';
	import { AlertStore } from '$lib/stores/alert';
	import { DialogStore } from '$lib/stores/dialog';
	import BoardCreateDialogContent from '$lib/components/board-save-dialog-content/board-create-dialog-content.svelte';
	import StatusCreateDialogContent from '$lib/components/status-save-dialog-content/status-create-dialog-content.svelte';
	import * as Scroll from '$lib/components/ui/scroll-area';
	import TicketCreateDialogContent from '$lib/components/ticket-save-dialog-content/ticket-create-dialog-content.svelte';
	import { cloneDeep, merge } from 'lodash';
	import {
		Subject,
		Subscription,
		catchError,
		combineLatest,
		concatMap,
		from,
		map,
		mergeAll,
		mergeMap,
		of,
		toArray,
		tap,
		bufferTime,
		buffer,
		timer,
		debounceTime
	} from 'rxjs';

	const { ScrollArea } = Scroll;

	let unsubscribes: Unsubscriber[] = [];
	let boardState = BoardStore.defaultState();
	let tempBoardState = BoardStore.defaultState();
	const ticketSubject = new Subject<Parameters<typeof TicketService.updateTicketsSortOrder>>();
	let tickets$: Subscription;
	const releaseBufferSubject = new Subject<void>();
	onMount(async () => {
		tickets$ = ticketSubject
			.pipe(
				mergeMap(([params, body]) =>
					from(TicketService.updateTicketsSortOrder(params, body)).pipe(
						map(({ data }) => {
							setTimeout(() => {
								releaseBufferSubject.next();
							}, 500);

							const { status_id } = params;

							return { status_id, data };
						}),
						catchError((error) => {
							AlertStore.error(error);

							return of();
						})
					)
				),
				buffer(releaseBufferSubject.pipe(debounceTime(1000)))
			)
			.subscribe((values) => {
				if (values.length === 0) return;
				console.log('values', values);
				BoardStore.update((state) => {
					const selected = cloneDeep(state.selected);
					for (const { status_id, data } of values) {
						const idx = selected.statuses.findIndex((status) => status.id === status_id);

						if (idx > -1) {
							selected.statuses[idx].tickets = data;
						}
					}

					state.selected = selected;

					return state;
				});
			});

		unsubscribes.push(
			authStore.subscribe((state) => {
				console.log('APP MOUNT', state);
				if (!state.user) {
					goto('/');
				}
			})
		);

		unsubscribes.push(
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
		unsubscribes.forEach((unsubscribe) => unsubscribe());
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
		ticketSubject.next([
			{
				board_id: boardState.selected.id,
				status_id: cid
			},
			{ tickets: boardState.selected.statuses[colIdx].tickets }
		]);
		// try {
		// 	const { data: tickets } = await TicketService.updateTicketsSortOrder(
		// 		{
		// 			board_id: boardState.selected.id,
		// 			status_id: cid
		// 		},
		// 		{ tickets: boardState.selected.statuses[colIdx].tickets }
		// 	);

		// 	BoardStore.update((state) => {
		// 		const targetStatus = boardState.selected.statuses[colIdx];

		// 		state.selected.statuses.forEach((status) => {
		// 			if (status.id === targetStatus.id) {
		// 				status.tickets = tickets;
		// 			}
		// 		});

		// 		return state;
		// 	});
		// } catch (error) {
		// 	AlertStore.error(error);

		// 	BoardStore.selectBoard(tempBoardState.selected);
		// }
	}

	enum Resource {
		Board = 'board',
		Status = 'status',
		Ticket = 'ticket'
	}
	enum Method {
		Create = 'create',
		Update = 'update'
	}

	const dialogState = {
		open: false,
		resrc_type: Resource.Board,
		method: Method.Create,
		data: null as any
	};

	function editTicket(ticket: TicketService.Ticket) {
		console.log('editTicket', ticket);
		dialogState.resrc_type = Resource.Ticket;
		dialogState.open = true;
	}
</script>

<section class="grid flex-1 grid-cols-12">
	<div class="relative col-span-2 px-2 bg-muted">
		<Button
			variant="ghost"
			class="absolute top-0 right-0 flex items-center justify-center h-auto p-1 m-1.5 hover:bg-opacity-10 hover:bg-accent-foreground"
			on:click={() => {
				DialogStore.create({ component: BoardCreateDialogContent });
			}}
		>
			<PlusOutline class="size-4" />
		</Button>
		<div class="flex flex-col w-full gap-2 mt-8">
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
					>
						<DotsHorizontalOutline class="size-4" />
					</Button>
				</Button>
			{/each}
		</div>
	</div>
	<div class="h-full col-span-10">
		<ScrollArea orientation="horizontal" class="has-[>div>div>div]:h-full">
			{#if boardState.selected}
				<div
					class="grid grid-flow-col gap-4 p-4 overflow-x-auto overflow-y-hidden"
					use:dndzone={{
						items: boardState.selected.statuses,
						flipDurationMs,
						type: 'columns',
						dropTargetStyle: {}
					}}
					on:consider={handleDndConsiderColumns}
					on:finalize={handleDndFinalizeColumns}
				>
					{#each boardState.selected.statuses as status (status.id)}
						<div animate:flip={{ duration: flipDurationMs }}>
							<ScrollArea
								orientation="vertical"
								class="h-[calc(100vh-(var(--header-height)+var(--footer-height))-2rem-1px)]"
							>
								<Card.Root class="p-2 w-80">
									<Card.Header class="px-2 py-2 pt-0">
										<Card.Title>
											<div class="flex items-center justify-between">
												<span class="text-base">{status.title}</span>
												<button
													class="p-1 ml-4 rounded cursor-pointer hover:text-accent-foreground hover:bg-accent"
													on:click={() => {
														DialogStore.create({
															component: TicketCreateDialogContent,
															params: { board_id: boardState?.selected?.id, status_id: status.id }
														});
													}}
												>
													<PlusOutline class="size-4" />
												</button>
											</div>
										</Card.Title>
									</Card.Header>

									<Card.Content class="relative px-0 py-0">
										<div
											class="absolute top-0 left-0 flex flex-col w-full gap-2"
											class:absolute={status.tickets.length === 0}
											class:h-20={status.tickets.length === 0}
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
													<TicketCard {ticket} edit={editTicket} />
												</div>
											{/each}
										</div>
										{#if status.tickets.length === 0}
											<button
												class="flex items-center justify-start w-full p-2 rounded hover:bg-accent"
												on:click={() => {
													DialogStore.create({
														component: TicketCreateDialogContent,
														params: { board_id: boardState?.selected?.id, status_id: status.id }
													});
												}}
											>
												<PlusOutline class="size-4" />
												<span class="ml-2">Add ticket</span>
											</button>
										{/if}
									</Card.Content>
								</Card.Root>
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
