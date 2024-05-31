<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index';
	import { Input } from '$lib/components/ui/input/index';
	import { Label } from '$lib/components/ui/label/index';
	import { TicketService } from '$lib/services/ticket-service';
	import { Button } from '$lib/components/ui/button/index';
	import { from } from 'rxjs';

	export let id = 0;
	export let data = {
		title: ''
	};
	export let boards: TicketService.Board[] = [];

	function handleSubmit() {
		if (!id) {
			from(TicketService.createBoard(data)).subscribe({
				next: ({ data }) => {
					boards = [...boards, data];
				},
				error: (error) => {
					console.error('createBoard', error);
				}
			});
		} else {
			console.log('update', data);
		}
	}
</script>

<Dialog.Content>
	<Dialog.Header>
		<Dialog.Title>{!id ? 'Create board' : 'Update Board'}</Dialog.Title>
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
		<Button type="submit" on:click={handleSubmit}>{!id ? 'Create' : 'Update'}</Button>
	</Dialog.Footer>
</Dialog.Content>
