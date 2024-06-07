<script lang="ts">
	import { DotsHorizontalOutline, PlusOutline } from 'flowbite-svelte-icons';
	import { flip } from 'svelte/animate';
	import { dndzone } from 'svelte-dnd-action';
	import Button from '$lib/components/ui/button/button.svelte';
	import TicketCard from '$lib/components/ticket-card/ticket-card.svelte';
	import type { Unsubscriber } from 'svelte/motion';
	import { onDestroy, onMount } from 'svelte';
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
		debounceTime,
		switchMap,
		takeUntil
	} from 'rxjs';
	import { AuthStore } from '$lib/stores/auth';
	import StatusCard from '$lib/components/status-card/status-card.svelte';
	import * as Select from '$lib/components/ui/select/index.js';

	const { ScrollArea } = Scroll;

	let unsubscribers: Unsubscriber[] = [];
	let localState = {
		...BoardStore.defaultState(),
		status_name: 'All status',
		status_index: -1
	};
	const statusSubject = new Subject<{ id: number; status_ids: number[] }>();
	const ticketSubject = new Subject<{
		id: number;
		ticket_ids: number[];
	}>();
	let status$: Subscription;
	let ticket$: Subscription;

	const releaseBufferSubject = new Subject<void>();
	onMount(async () => {
		unsubscribers.push(AuthStore.Use());
		status$ = statusSubject
			.pipe(
				debounceTime(1000),
				switchMap(({ id, status_ids }) =>
					from(TicketService.updateStatusesSortOrder({ board_id: id }, { status_ids })).pipe(
						map(({ data }) => {
							return data;
						}),
						catchError((error) => {
							AlertStore.error(error);
							fetchBoardFullDetail(localState.selected);

							return of();
						}),
						takeUntil(statusSubject)
					)
				)
			)
			.subscribe();

		ticket$ = ticketSubject
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
							{ board_id: localState.selected.id },
							{ statuses: values }
						)
					).pipe(
						map(({ data }) => {
							return data;
						}),
						catchError((error) => {
							AlertStore.error(error);
							fetchBoardFullDetail(localState.selected);

							return of();
						})
					)
				)
			)
			.subscribe();

		unsubscribers.push(
			BoardStore.subscribe((state) => {
				console.log('board state change', state);

				localState = {
					...localState,
					...cloneDeep(state)
				};
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
		status$.unsubscribe();
		ticket$.unsubscribe();
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
		if (!localState.selected) return;
		localState.selected.statuses = [...e.detail.items];
	}
	async function handleDndFinalizeColumns(e: ColumnEvent) {
		if (!localState.selected) return;
		localState.selected.statuses = [...e.detail.items];

		statusSubject.next({
			id: localState.selected.id,
			status_ids: localState.selected.statuses.map((s) => s.id)
		});
	}

	type CardEvent = CustomEvent & { detail: { items: TicketService.Ticket[] } };
	function handleDndConsiderCards(cid: number, e: CardEvent) {
		if (!localState.selected) return;
		const colIdx = localState.selected.statuses?.findIndex((c) => c.id === cid);
		localState.selected.statuses[colIdx].tickets = e.detail.items;
		localState.selected.statuses = [...localState.selected.statuses];
	}
	async function handleDndFinalizeCards(cid: number, e: CardEvent) {
		if (!localState.selected) return;
		const colIdx = localState.selected.statuses.findIndex((c) => c.id === cid);
		localState.selected.statuses[colIdx].tickets = e.detail.items;
		localState.selected.statuses = [...localState.selected.statuses];

		console.log('handleDndFinalizeCards');
		ticketSubject.next({
			id: cid,
			ticket_ids: localState.selected.statuses[colIdx].tickets.map((t) => t.id)
		});
	}

	function onSortTicketsInStatus(status: TicketService.Status) {
		let idx = localState.selected.statuses.findIndex((s) => s.id === status.id);

		if (idx === -1) return;

		localState.selected.statuses[idx] = status;

		ticketSubject.next({
			id: status.id,
			ticket_ids: status.tickets.map((t) => t.id)
		});
	}
</script>

<section class="grid flex-1 grid-cols-12">
	<div class="relative col-span-2 pl-2 bg-muted">
		<div class="flex justify-end w-full h-10">
			<Button
				variant="ghost"
				class="flex items-center justify-center h-auto p-1 my-auto mr-2 hover:bg-opacity-10 hover:bg-accent-foreground"
				on:click={() => {
					DialogStore.create({ component: BoardSaveDialogContent });
				}}
			>
				<PlusOutline class="size-4" />
			</Button>
		</div>
		<ScrollArea
			class=" mr-2 w-full h-[calc(100vh-(var(--header-height)+var(--footer-height))-2.5rem-1px)] [&>[data-melt-scroll-area-thumb]]:bg-red-400"
			orientation="vertical"
			scrollbarYClasses="dark:[&>[data-melt-scroll-area-thumb]]:bg-primary-foreground"
		>
			<div class="flex flex-col w-full gap-2">
				{#each localState.boards as board (board.id)}
					<Button
						variant="ghost"
						class="justify-between py-2 hover:bg-opacity-10 hover:bg-accent-foreground group/sidemenu {localState.selected &&
						localState.selected.id === board.id
							? 'bg-accent-foreground bg-opacity-10 text-accent-foreground'
							: ''}"
						on:click={() => {
							localState.status_index = -1;
							localState.status_name = 'All status';
							fetchBoardFullDetail(board);
						}}
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
		<div class="flex items-center justify-start w-full h-10">
			<Select.Root
				selected={{
					label: localState.status_name,
					value: localState.status_index
				}}
				onSelectedChange={(e) => {
					localState.status_name = e.label;
					localState.status_index = e.value;
				}}
			>
				<Select.Trigger class="w-[180px] max-h-8 ml-2">
					<Select.Value placeholder="All status" />
				</Select.Trigger>
				<Select.Content>
					<Select.Group>
						<Select.Item value={-1} label="All status" class="cursor-pointer">
							All status
						</Select.Item>
						{#each localState.selected.statuses as status, i (status.id)}
							<Select.Item value={i} label={status.title} class="cursor-pointer">
								{status.title}
							</Select.Item>
						{/each}
					</Select.Group>
				</Select.Content>
				<Select.Input name="filterStatus" />
			</Select.Root>
		</div>

		<ScrollArea orientation={localState.status_index > -1 ? 'vertical' : 'horizontal'}>
			{#if localState.selected && localState.status_index == -1}
				<div
					class="flex justify-start gap-4 pl-4 overflow-x-auto overflow-y-hidden"
					use:dndzone={{
						items: localState.selected.statuses,
						flipDurationMs,
						type: 'columns',
						dropTargetStyle: {}
					}}
					on:consider={handleDndConsiderColumns}
					on:finalize={handleDndFinalizeColumns}
				>
					{#each localState.selected.statuses as status, i (status)}
						<div animate:flip={{ duration: flipDurationMs }}>
							<ScrollArea
								orientation="vertical"
								class="h-[calc(100vh-(var(--header-height)+var(--footer-height))-2.5rem-1px)]"
							>
								<StatusCard {status} {onSortTicketsInStatus}>
									{#if status.tickets.length === 0}
										<Button
											variant="ghost"
											class="flex items-center justify-start w-full h-auto p-1 rounded hover:bg-accent"
											on:click={() => {
												DialogStore.create({
													component: TicketSaveDialogContent,
													params: { board_id: localState?.selected?.id, status_id: status.id }
												});
											}}
										>
											<PlusOutline class="size-4" />
											<span class="ml-2">Add ticket</span>
										</Button>
									{/if}
									<div
										class="absolute top-0 left-0 flex flex-col w-full gap-2"
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
									params: { board_id: localState?.selected?.id }
								});
							}}
						>
							<PlusOutline class="size-4" />
							<span class="ml-2">Add status</span>
						</button>
					</div>
				</div>
			{/if}
			{#if localState.status_index > -1}
				<div class="flex justify-start gap-4 px-4 overflow-x-hidden overflow-y-auto">
					<div
						class="absolute top-0 left-0 flex flex-col w-full gap-2"
						class:absolute={localState.selected.statuses[localState.status_index].tickets.length ===
							0}
						class:h-32={localState.selected.statuses[localState.status_index].tickets.length === 0}
						use:dndzone={{
							items: localState.selected.statuses[localState.status_index].tickets,
							flipDurationMs,
							dropTargetStyle: {}
						}}
						on:consider={(e) =>
							handleDndConsiderCards(localState.selected.statuses[localState.status_index].id, e)}
						on:finalize={(e) =>
							handleDndFinalizeCards(localState.selected.statuses[localState.status_index].id, e)}
					>
						{#each localState.selected.statuses[localState.status_index].tickets as ticket (ticket.id)}
							<div tabindex={ticket.id} role="button" animate:flip={{ duration: flipDurationMs }}>
								<TicketCard
									{ticket}
									board_id={localState.selected.statuses[localState.status_index].board_id}
								/>
							</div>
						{/each}
					</div>
				</div>
			{/if}
		</ScrollArea>
	</div>
</section>
