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
	import { TicketService } from '$lib/services/tircket-service';
	import { from } from 'rxjs';

	let unsubscribe: Unsubscriber;
	onMount(() => {
		unsubscribe = authStore.subscribe((state) => {
			console.log('APP MOUNT', state);
			if (!state.user) {
				goto('/app');
			}
		});

		from(TicketService.getBoards()).subscribe({
			next: ({ data }) => {
				boards.push(...data);
			},
			error: (error) => {
				console.error('getBoards', error);
			}
		});
	});

	onDestroy(() => {
		unsubscribe();
	});

	const flipDurationMs = 200;
	const boards: TicketService.Board[] = [];

	let index = 0;

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
		method: Method.Create
	};

	function editTicket(ticket: TicketService.Ticket) {
		console.log('editTicket', ticket);
		dialogState.resrc_type = Resource.Ticket;
		dialogState.open = true;
	}
</script>

<section class="h-[90vh] w-full">
	<Dialog.Root open={dialogState.open}>
		<div class="h-full mt-2">
			<BoardTabs.Root>
				<BoardTabs.List>
					{#each boards as board (board.id)}
						<BoardTabs.Trigger value={`${boards[index].id}`}>
							<span>{boards[index].title}</span>
						</BoardTabs.Trigger>
					{/each}
					<div
						class="flex items-center justify-center p-1 m-3 rounded cursor-pointer hover:text-accent-foreground hover:bg-accent"
					>
						<PlusOutline class="size-4" />
					</div>
				</BoardTabs.List>
				{#each boards as board (board.id)}
					<BoardTabs.Content value={`${boards[index].id}`}>
						<div
							class="flex justify-start h-full gap-4 p-4 overflow-x-auto overflow-y-hidden"
							use:dndzone={{
								items: boards[index].statuses,
								flipDurationMs,
								type: 'columns',
								dropTargetStyle: {}
							}}
							on:consider={handleDndConsiderColumns}
							on:finalize={handleDndFinalizeColumns}
						>
							{#each boards[index].statuses as status (status.id)}
								<div class="relative" animate:flip={{ duration: flipDurationMs }}>
									<Card.Root class="px-2 w-80">
										<Card.Header>
											<Card.Title>{status.title}</Card.Title>
										</Card.Header>
										<Card.Content class="px-0">
											<div
												class="flex flex-col gap-2 min-h-32"
												use:dndzone={{ items: status.tickets, flipDurationMs, dropTargetStyle: {} }}
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
					</BoardTabs.Content>
				{/each}
			</BoardTabs.Root>
		</div>
		<Dialog.Content>
			<Dialog.Header>
				<Dialog.Title>Are you sure absolutely sure?</Dialog.Title>
				<Dialog.Description>
					This action cannot be undone. This will permanently delete your account and remove your
					data from our servers.
				</Dialog.Description>
			</Dialog.Header>
		</Dialog.Content>
	</Dialog.Root>
</section>
