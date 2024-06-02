<script lang="ts">
	import * as Card from '$lib/components/ui/card/index';
	import type { TicketService } from '$lib/services/ticket-service';
	import { DotsHorizontalOutline } from 'flowbite-svelte-icons';
	import Button from '$lib/components/ui/button/button.svelte';
	export let ticket = {
		id: 0,
		title: '',
		description: ''
	} as TicketService.Ticket;
	export let edit = (ticket: TicketService.Ticket) => {};
</script>

<Card.Root class="p-3">
	<Card.Header class="p-0 group">
		<Card.Title>
			<div class="relative flex justify-between">
				<span class="text-sm">{ticket.title}</span>
				<Button
					variant="ghost"
					class="flex items-center justify-center invisible h-auto p-1 hover:bg-opacity-10 hover:bg-accent-foreground group-hover:visible"
				>
					<DotsHorizontalOutline class="size-4" />
				</Button>
			</div>
		</Card.Title>
	</Card.Header>
	<Card.CardContent class="p-0 mt-3 text-sm line-clamp-6">
		{ticket.description}
	</Card.CardContent>
</Card.Root>
