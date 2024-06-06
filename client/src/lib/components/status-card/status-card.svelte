<script lang="ts">
	import * as Card from '$lib/components/ui/card/index';
	import Button from '$lib/components/ui/button/button.svelte';
	import { DialogStore } from '$lib/stores/dialog';
	import type { Ticket } from 'lucide-svelte';
	import type { TicketService } from '$lib/services/ticket-service';
	import StatusSaveDialogContent from '../status-save-dialog-content/status-save-dialog-content.svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index';
	import { DotsHorizontalOutline, PlusOutline, SortOutline } from 'flowbite-svelte-icons';
	import TicketSaveDialogContent from '../ticket-save-dialog-content/ticket-save-dialog-content.svelte';

	export let status: TicketService.Status;

	function sortTicketsInstatus({
		status,
		column,
		direction
	}: {
		status: TicketService.Status;
		column: 'created_at' | 'updated_at';
		direction: 'asc' | 'desc';
	}) {
		status.tickets = status.tickets.sort((a, b) => {
			if (direction === 'asc') {
				return new Date(a[column]).getTime() - new Date(b[column]).getTime();
			} else {
				return new Date(b[column]).getTime() - new Date(a[column]).getTime();
			}
		});

		onSortTicketsInStatus(status);
	}
	export let onSortTicketsInStatus: (status: TicketService.Status) => void;
</script>

<Card.Root class="p-2 w-80">
	<Card.Header class="px-2 py-2 pt-0 group">
		<Card.Title>
			<div class="flex items-center justify-between">
				<span class="text-base">{status.title}</span>
				<div class="flex">
					<Button
						variant="ghost"
						class="flex items-center justify-center invisible h-auto p-1 ml-4 rounded cursor-pointer group-hover:visible hover:text-accent-foreground hover:bg-accent"
						on:click={() => {
							DialogStore.create({
								component: StatusSaveDialogContent,
								params: { model: status }
							});
						}}
					>
						<DotsHorizontalOutline class="size-4" />
					</Button>
					<DropdownMenu.Root>
						<DropdownMenu.Trigger asChild let:builder>
							<Button
								builders={[builder]}
								variant="ghost"
								class="flex items-center justify-center invisible h-auto p-1 ml-2 rounded cursor-pointer hover:text-accent-foreground hover:bg-accent group-hover:visible"
							>
								<SortOutline class="size-4" />
							</Button>
						</DropdownMenu.Trigger>
						<DropdownMenu.Content class="w-10">
							<DropdownMenu.Label>Sort tickets by</DropdownMenu.Label>
							<DropdownMenu.Separator />
							<DropdownMenu.Group>
								<DropdownMenu.Item>
									<button
										on:click={() => {
											sortTicketsInstatus({
												status,
												column: 'created_at',
												direction: 'desc'
											});
										}}>Lastest created</button
									>
								</DropdownMenu.Item>
								<DropdownMenu.Item>
									<button
										on:click={() => {
											sortTicketsInstatus({
												status,
												column: 'created_at',
												direction: 'asc'
											});
										}}>Oldest created</button
									>
								</DropdownMenu.Item>
								<DropdownMenu.Item>
									<button
										on:click={() => {
											sortTicketsInstatus({
												status,
												column: 'updated_at',
												direction: 'desc'
											});
										}}>Lastest updated</button
									>
								</DropdownMenu.Item>
								<DropdownMenu.Item>
									<button
										on:click={() => {
											sortTicketsInstatus({
												status,
												column: 'updated_at',
												direction: 'asc'
											});
										}}>Oldest updated</button
									>
								</DropdownMenu.Item>
							</DropdownMenu.Group>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
					<Button
						variant="ghost"
						class="flex items-center justify-center h-auto p-1 ml-2 rounded cursor-pointer hover:text-accent-foreground hover:bg-accent"
						on:click={() => {
							DialogStore.create({
								component: TicketSaveDialogContent,
								params: { board_id: status.board_id, status_id: status.id }
							});
						}}
					>
						<PlusOutline class="size-4" />
					</Button>
				</div>
			</div>
		</Card.Title>
	</Card.Header>
	<Card.Content class="relative px-0 py-0">
		<slot />
	</Card.Content>
</Card.Root>
