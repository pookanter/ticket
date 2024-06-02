<script lang="ts">
	import * as BoardTabs from '$lib/components/board-tabs/index';
	import { PlusOutline } from 'flowbite-svelte-icons';
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
	import BoardSaveDialogContent from '$lib/components/board-save-dialog-content/board-create-dialog-content.svelte';
	import { BoardStore } from '$lib/stores/board';
	import { AlertStore } from '$lib/stores/alert';
	import { DialogStore } from '$lib/stores/dialog';
	import BoardCreateDialogContent from '$lib/components/board-save-dialog-content/board-create-dialog-content.svelte';
	import StatusCreateDialogContent from '$lib/components/status-save-dialog-content/status-create-dialog-content.svelte';
	import { ScrollArea } from '$lib/components/ui/scroll-area';

	let unsubscribes: Unsubscriber[] = [];
	let boards: TicketService.Board[] = [];
	let index = 0;
	let boardState = BoardStore.defaultState();
	onMount(() => {
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
				if (state.initializing) {
					try {
						const { data: boards } = await TicketService.getBoards();
						BoardStore.update((state) => {
							state.boards = boards;

							if (boards.length > 0) {
								state.selected = boards[0];
							}

							return state;
						});
					} catch ({ error, message }: any) {
						AlertStore.create({
							title: 'Error',
							message: error ? error.message : message || 'An error occurred'
						});
					}

					BoardStore.update((state) => {
						state.initializing = false;
						return state;
					});
					return;
				}

				console.log('board state change', state.boards);

				boardState = { ...state };
			})
		);
	});

	onDestroy(() => {
		unsubscribes.forEach((unsubscribe) => unsubscribe());
		BoardStore.update((state) => {
			state.initializing = true;
			return state;
		});
	});

	const flipDurationMs = 200;

	type ColumnEvent = CustomEvent & { detail: { items: TicketService.Status[] } };
	function handleDndConsiderColumns(e: ColumnEvent) {
		console.log('boards[index] before', boards[index].statuses);
		boards[index].statuses = e.detail.items;

		console.log('boards[index] change', boards[index].statuses);
	}
	function handleDndFinalizeColumns(e: ColumnEvent) {
		boards[index].statuses = e.detail.items;
	}

	type CardEvent = CustomEvent & { detail: { items: TicketService.Ticket[] } };
	function handleDndConsiderCards(cid: number, e: CardEvent) {
		console.log('handleDndConsiderCards', cid, e.detail.items);
		const colIdx = boards[index].statuses?.findIndex((c) => c.id === cid);
		boards[index].statuses[colIdx].tickets = e.detail.items;
		boards[index].statuses = [...boards[index].statuses];
	}
	function handleDndFinalizeCards(cid: number, e: CardEvent) {
		const colIdx = boards[index].statuses.findIndex((c) => c.id === cid);
		boards[index].statuses[colIdx].tickets = e.detail.items;
		boards[index].statuses = [...boards[index].statuses];
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

<section class="h-[90vh] w-full">
	<BoardTabs.Root class="mt-2 size-full">
		<BoardTabs.List>
			{#each boardState.boards as board (board.id)}
				<BoardTabs.Trigger
					value={board}
					clickupdate={(value) => {}}
					on:click={() => {
						BoardStore.update((state) => {
							state.selected = board;
							console.log('selected', state.selected);
							return state;
						});
					}}
				>
					<span>{board.title}</span>
				</BoardTabs.Trigger>
			{/each}
			<button
				on:click={() => {
					DialogStore.create({ component: BoardCreateDialogContent });
				}}
				class="flex items-center justify-center p-1 m-3 rounded cursor-pointer hover:text-accent-foreground hover:bg-accent"
			>
				<PlusOutline class="size-4" />
			</button>
		</BoardTabs.List>
		{#each boardState.boards as board (board.id)}
			<BoardTabs.Content value={`${board.id}`}>
				<ScrollArea class="size-full" orientation="horizontal">
					<div
						class="flex justify-start h-full gap-4 p-4 overflow-x-auto overflow-y-hidden"
						use:dndzone={{
							items: board.statuses,
							flipDurationMs,
							type: 'columns',
							dropTargetStyle: {}
						}}
						on:consider={handleDndConsiderColumns}
						on:finalize={handleDndFinalizeColumns}
					>
						{#each board.statuses as status (status.id)}
							<div class="relative" animate:flip={{ duration: flipDurationMs }}>
								<Card.Root class="px-2 w-80">
									<Card.Header>
										<Card.Title>
											<div class="flex items-center justify-between ">
												<span>{status.title}</span>
												<button
													class="p-1 ml-4 rounded cursor-pointer hover:text-accent-foreground hover:bg-accent"
												>
													<PlusOutline class="size-4" />
												</button>
											</div>
										</Card.Title>
									</Card.Header>
									<Card.Content class="px-0">
										<div
											class="flex flex-col gap-2"
											use:dndzone={{
												items: status.tickets,
												flipDurationMs,
												dropTargetStyle: {}
											}}
											on:consider={(e) => handleDndConsiderCards(status.id, e)}
											on:finalize={(e) => handleDndFinalizeCards(status.id, e)}
										>
											{#each status.tickets as ticket (ticket.id)}
												<div animate:flip={{ duration: flipDurationMs }}>
													<TicketCard {ticket} edit={editTicket} />
												</div>
											{/each}
										</div>
										{#if status.tickets.length === 0}
											<button
												class="flex items-center justify-start w-full p-2 rounded hover:bg-accent"
												on:click={() => {
													DialogStore.create({
														component: StatusCreateDialogContent,
														params: { board_id: board.id }
													});
												}}
											>
												<PlusOutline class="size-4" />
												<span class="ml-2">Add ticket</span>
											</button>
										{/if}
									</Card.Content>
								</Card.Root>
							</div>
						{/each}
						<div class="block">
							<button
								class="flex items-center justify-start p-2 mt-4 rounded min-w-80 hover:bg-accent"
								on:click={() => {
									DialogStore.create({
										component: StatusCreateDialogContent,
										params: { board_id: board.id }
									});
								}}
							>
								<PlusOutline class="size-4" />
								<span class="ml-2">Add status</span>
							</button>
						</div>
					</div>
				</ScrollArea>
			</BoardTabs.Content>
		{/each}
	</BoardTabs.Root>
</section>
