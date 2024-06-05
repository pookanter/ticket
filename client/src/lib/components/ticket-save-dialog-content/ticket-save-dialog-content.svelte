<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index';
	import { Input } from '$lib/components/ui/input/index';
	import { Label } from '$lib/components/ui/label/index';
	import { Button } from '$lib/components/ui/button/index';
	import { AlertStore } from '$lib/stores/alert';
	import { DialogStore } from '$lib/stores/dialog';
	import { TicketService } from '$lib/services/ticket-service';
	import { BoardStore } from '$lib/stores/board';
	import Textarea from '$lib/components/ui/textarea/textarea.svelte';
	import { onMount } from 'svelte';

	export let board_id: number;
	export let status_id: number;
	export let model = {
		id: 0,
		status_id: 0,
		title: '',
		description: '',
		contact: ''
	};

	onMount(() => {
		data.title = model.title;
		data.description = model.description;
		data.contact = model.contact;
	});

	let data: Parameters<typeof TicketService.createTicket>[1] = {
		title: '',
		description: '',
		contact: ''
	};

	async function handleSubmit() {
		if (model.id > 0) {
			try {
				const { data: ticket } = await TicketService.updateTicketPartial(
					{ board_id, status_id: model.status_id, ticket_id: model.id },
					data
				);

				BoardStore.updateTicket({ ticket });
				DialogStore.close();
			} catch (error: any) {
				AlertStore.error(error);
			}
			return;
		}

		try {
			const { data: ticket } = await TicketService.createTicket({ board_id, status_id }, data);

			BoardStore.addTicket({ board_id, ticket });
			DialogStore.close();
		} catch (error: any) {
			AlertStore.error(error);
		}
	}
</script>

<Dialog.Content>
	<Dialog.Header>
		<Dialog.Title>{model.id > 0 ? 'Edit' : 'Add'} ticket</Dialog.Title>
	</Dialog.Header>
	<div class="grid gap-4 py-4">
		<div class="grid gap-4 py-4">
			<div class="grid items-center grid-cols-4 gap-4">
				<Label for="name" class="text-right">Title</Label>
				<Input id="name" bind:value={data.title} class="col-span-3" />
			</div>
			<div class="grid items-center grid-cols-4 gap-4">
				<Label for="name" class="text-right">Description</Label>
				<Textarea bind:value={data.description} class="col-span-3" />
			</div>
			<div class="grid items-center grid-cols-4 gap-4">
				<Label for="name" class="text-right">Contact</Label>
				<Input id="name" bind:value={data.contact} class="col-span-3" />
			</div>
		</div>
	</div>
	<Dialog.Footer>
		<Button type="submit" on:click={handleSubmit}>{model.id > 0 ? 'Edit' : 'Add'}</Button>
	</Dialog.Footer>
</Dialog.Content>
