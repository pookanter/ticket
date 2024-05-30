<script lang="ts">
	import * as Card from '$lib/components/ui/card/index';
	import * as Dialog from '$lib/components/ui/dialog/index';
	import { EditOutline, PlusOutline } from 'flowbite-svelte-icons';
	interface Ticket {
		id: number;
		title: string;
		description: string;
	}

	export let ticket: Ticket = {
		id: 0,
		title: '',
		description: ''
	};
	export let edit = (ticket: Ticket) => {};
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
