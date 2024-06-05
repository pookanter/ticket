<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index';
	import { Input } from '$lib/components/ui/input/index';
	import { Label } from '$lib/components/ui/label/index';
	import { TicketService } from '$lib/services/ticket-service';
	import { Button } from '$lib/components/ui/button/index';
	import { BoardStore } from '$lib/stores/board';
	import { AlertStore } from '$lib/stores/alert';
	import { DialogStore } from '$lib/stores/dialog';
	import { onMount } from 'svelte';

	export let board_id: number;
	export let model = {
		id: 0,
		board_id: 0,
		title: ''
	};

	onMount(() => {
		console.log(model);
		data.title = model.title;
	});

	let data = {
		title: ''
	};

	async function handleSubmit() {
		if (model.id > 0) {
			try {
				const { data: status } = await TicketService.updateStatusPartial(
					{ board_id: model.board_id, status_id: model.id },
					data
				);

				BoardStore.updateStatus({ status });
				DialogStore.close();
			} catch (error: any) {
				AlertStore.error(error);
			}
			return;
		}

		try {
			const { data: status } = await TicketService.createStatus({ board_id }, data);

			BoardStore.addStatus({ status });
			DialogStore.close();
		} catch (error: any) {
			AlertStore.error(error);
		}
	}
</script>

<Dialog.Content>
	<Dialog.Header>
		<Dialog.Title>{model.id > 0 ? 'Edit' : 'Add'} status</Dialog.Title>
	</Dialog.Header>
	<div class="grid gap-4 py-4">
		<div class="grid gap-4 py-4">
			<div class="grid items-center grid-cols-4 gap-4">
				<Label for="name" class="text-right">Title</Label>
				<Input id="name" bind:value={data.title} class="col-span-3" />
			</div>
		</div>
	</div>
	<Dialog.Footer>
		<Button type="submit" on:click={handleSubmit}>{model.id > 0 ? 'Edit' : 'Add'}</Button>
	</Dialog.Footer>
</Dialog.Content>
