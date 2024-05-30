<script lang="ts">
	import { AuthenService } from '$lib/services/authen-service';
	import { auth, me } from '$lib/store/authen';
	import { from } from 'rxjs';
	import { onMount } from 'svelte';
	import * as BoardTabs from '$lib/components/board-tabs/index';
	import { FolderOpenOutline, PlusOutline } from 'flowbite-svelte-icons';
	import { v4 as uuidv4 } from 'uuid';
	import * as Card from '$lib/components/ui/card/index';
	import { flip } from 'svelte/animate';
	import { dndzone } from 'svelte-dnd-action';
	import Button from '$lib/components/ui/button/button.svelte';
	interface Ticket {
		id: number;
		title: string;
		description: string;
	}

	interface Status {
		id: number;
		title: string;
		tickets: Ticket[];
	}

	interface Board {
		id: number;
		title: string;
		statuses: Status[];
	}

	const flipDurationMs = 200;
	const boards = [
		{
			id: 1,
			title: 'Board 1',
			statuses: [
				{
					id: 1,
					title: 'Status 1',
					tickets: [
						{
							id: 1,
							title: 'Task 1',
							description: 'This is a mock task.'
						},
						{
							id: 2,
							title: 'Task 2',
							description: 'This is a mock task.'
						}
					]
				},
				{
					id: 2,
					title: 'Status 2',
					tickets: [
						{
							id: 3,
							title: 'Task 3',
							description: 'This is a mock task.'
						},
						{
							id: 4,
							title: 'Task 4',
							description: 'This is a mock task.'
						}
					]
				}
			]
		}
		// {
		// 	id: 2,
		// 	title: 'Board 2',
		// 	statuses: []
		// }
	] as Board[];

	let index = 0;

	type ColumnEvent = CustomEvent & { detail: { items: Status[] } };
	function handleDndConsiderColumns(e: ColumnEvent) {
		console.log('boards[index] before', boards[index].statuses);
		boards[index].statuses = e.detail.items;

		console.log('boards[index] change', boards[index].statuses);
	}
	function handleDndFinalizeColumns(e: ColumnEvent) {
		boards[index].statuses = e.detail.items;
	}

	type CardEvent = CustomEvent & { detail: { items: Ticket[] } };
	function handleDndConsiderCards(cid: number, e: CardEvent) {
		console.log('handleDndConsiderCards', cid, e.detail.items);
		const colIdx = boards[index].statuses.findIndex((c) => c.id === cid);
		boards[index].statuses[colIdx].tickets = e.detail.items;
		boards[index].statuses = [...boards[index].statuses];
	}
	function handleDndFinalizeCards(cid: number, e: CardEvent) {
		const colIdx = boards[index].statuses.findIndex((c) => c.id === cid);
		boards[index].statuses[colIdx].tickets = e.detail.items;
		boards[index].statuses = [...boards[index].statuses];
	}
</script>

<section class="h-[90vh] w-full">
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
													<Card.Root>
														<Card.Header>
															<Card.Title>{ticket.title}</Card.Title>
														</Card.Header>
														<Card.Content>{ticket.id}</Card.Content>
													</Card.Root>
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
</section>
