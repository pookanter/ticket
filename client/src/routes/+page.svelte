<script lang="ts">
	import * as BoardTabs from '$lib/components/board-tabs/index';
	import { PlusOutline } from 'flowbite-svelte-icons';
	import * as Card from '$lib/components/ui/card/index';
	import { flip } from 'svelte/animate';
	import { dndzone } from 'svelte-dnd-action';
	import Button from '$lib/components/ui/button/button.svelte';
	import TicketCard from '$lib/components/ticket-card/ticket-card.svelte';
	import * as Dialog from '$lib/components/ui/dialog/index';
	import type { Unsubscriber } from 'svelte/motion';
	import { onDestroy, onMount } from 'svelte';
	import authStore from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { TicketService } from '$lib/services/ticket-service';
	import { from } from 'rxjs';
	import BoardSaveDialogContent from '$lib/components/board-save-dialog/board-save-dialog-content.svelte';
	import { BoardStore } from '$lib/stores/board';
	import { AlertStore } from '$lib/stores/alert';

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

	const boardDialog = {
		method: Method.Create,
		data: {
			id: 0,
			title: ''
		}
	};

	let dialogOpen = false;
	let targetResource: Resource = Resource.Board;
	function onDialogOpenChange(open: boolean) {
		dialogOpen = open;
	}
	function openBoardDialog({
		method,
		value
	}: {
		method: Method;
		value?: TicketService.Board | null;
	}) {
		console.log('dialogOpen', dialogOpen);
		console.log('openBoardDialog', method, value);
		if (method == Method.Create) {
			boardDialog.method = Method.Create;
			dialogOpen = true;
		}

		if (method == Method.Update && value) {
			boardDialog.method = Method.Update;
			boardDialog.data = { ...value };
			dialogOpen = true;
		}

		targetResource = Resource.Board;
	}

	function editTicket(ticket: TicketService.Ticket) {
		console.log('editTicket', ticket);
		dialogState.resrc_type = Resource.Ticket;
		dialogState.open = true;
	}
</script>

<section class="h-[90vh] w-full">
	<Dialog.Root open={dialogOpen} onOpenChange={onDialogOpenChange}>
		<div class="h-full mt-2">
			<BoardTabs.Root>
				<BoardTabs.List>
					{#each boardState.boards as board (board.id)}
						<BoardTabs.Trigger
							value={board}
							clickupdate={(value) => openBoardDialog({ method: Method.Update, value })}
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
						on:click={() => openBoardDialog({ method: Method.Create })}
						class="flex items-center justify-center p-1 m-3 rounded cursor-pointer hover:text-accent-foreground hover:bg-accent"
					>
						<PlusOutline class="size-4" />
					</button>
				</BoardTabs.List>
				{#each boardState.boards as board (board.id)}
					<BoardTabs.Content value={`${board.id}`}>
						{#if boardState.selected}
							<div
								class="flex justify-start h-full gap-4 p-4 overflow-x-auto overflow-y-hidden"
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
									<div class="relative" animate:flip={{ duration: flipDurationMs }}>
										<Card.Root class="px-2 w-80">
											<Card.Header>
												<Card.Title>{status.title}</Card.Title>
											</Card.Header>
											<Card.Content class="px-0">
												<div
													class="flex flex-col gap-2 min-h-32"
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
												<Button variant="outline" class="w-full mt-2">
													<PlusOutline class="size-4" />
												</Button>
											</Card.Content>
										</Card.Root>
									</div>
								{/each}
							</div>
						{/if}
					</BoardTabs.Content>
				{/each}
			</BoardTabs.Root>
		</div>
		{#if targetResource === Resource.Board}
			<BoardSaveDialogContent id={boardDialog.data.id} data={boardDialog.data} {boards} />
		{/if}
	</Dialog.Root>
</section>
