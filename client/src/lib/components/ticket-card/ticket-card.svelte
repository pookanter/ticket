<script lang="ts">
	import * as Card from '$lib/components/ui/card/index';
	import type { TicketService } from '$lib/services/ticket-service';
	import { ClockOutline, DotsHorizontalOutline } from 'flowbite-svelte-icons';
	import Button from '$lib/components/ui/button/button.svelte';
	import { DialogStore } from '$lib/stores/dialog';
	import TicketSaveDialogContent from '../ticket-save-dialog-content/ticket-save-dialog-content.svelte';
	import dayjs from 'dayjs/esm';
	import * as Tooltip from '$lib/components/ui/tooltip/index';
	export let board_id: number;
	export let ticket = {
		id: 0,
		title: '',
		description: ''
	} as TicketService.Ticket;
</script>

<Card.Root class="p-3">
	<Card.Header class="p-0 group">
		<Card.Title>
			<div class="relative flex justify-between">
				<span class="text-sm">{ticket.title}</span>
				<Button
					variant="ghost"
					class="flex items-center justify-center invisible h-auto p-1 hover:bg-opacity-10 hover:bg-accent-foreground group-hover:visible"
					on:click={() =>
						DialogStore.create({
							component: TicketSaveDialogContent,
							params: { model: ticket, board_id }
						})}
				>
					<DotsHorizontalOutline class="size-4" />
				</Button>
			</div>
		</Card.Title>
	</Card.Header>
	<Card.CardContent class="p-0 mt-3 text-sm line-clamp-6">
		{ticket.description}
		<div class="mt-3">
			<Tooltip.Root>
				<Tooltip.Trigger>
					<div class="flex gap-1 mt-3">
						<ClockOutline class="text-gray-500 size-4 dark:text-gray-300" />
						<span class="text-xs text-gray-500 dark:text-gray-300"
							>{dayjs(ticket.created_at).format('MMM DD')}</span
						>
					</div>
				</Tooltip.Trigger>
				<Tooltip.Content>
					{#if ticket.updated_at}
						lastest updated at {dayjs(ticket.updated_at).format('MMM DD, YYYY')}
					{:else}
						created at {dayjs(ticket.created_at).format('MMM DD, YYYY')}
					{/if}
				</Tooltip.Content>
			</Tooltip.Root>
		</div>
	</Card.CardContent>
</Card.Root>
