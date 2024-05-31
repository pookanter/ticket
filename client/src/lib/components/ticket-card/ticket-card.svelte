<script lang="ts">
	import * as Card from '$lib/components/ui/card/index';
	import * as Dialog from '$lib/components/ui/dialog/index';
	import type { TicketService } from '$lib/services/ticket-service';
	import { EditOutline, PlusOutline } from 'flowbite-svelte-icons';
	export let ticket = {
		id: 0,
		title: '',
		description: ''
	} as TicketService.Ticket;
	export let edit = (ticket: TicketService.Ticket) => {};
</script>

<Card.Root>
	<Card.Header class="group">
		<Card.Title>
			<div class="relative flex start">
				<span>{ticket.title}</span>
				<button
					class="absolute top-0 right-0 flex items-center justify-center invisible p-1 -mt-4 -mr-4 rounded cursor-pointer group-hover:visible hover:text-accent-foreground hover:bg-accent"
					on:click={() => edit(ticket)}
				>
					<EditOutline class="size-4" />
				</button>
			</div>
		</Card.Title>
	</Card.Header>
</Card.Root>
