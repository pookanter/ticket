<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index';
	import { Input } from '$lib/components/ui/input/index';
	import { Label } from '$lib/components/ui/label/index';
	import { TicketService } from '$lib/services/ticket-service';
	import { Button } from '$lib/components/ui/button/index';
	import { BoardStore } from '$lib/stores/board';
	import { AlertStore } from '$lib/stores/alert';
	import { DialogStore } from '$lib/stores/dialog';
	import { cloneDeep } from 'lodash';

	export let board_id: number;

	let data = {
		title: ''
	};

	async function handleSubmit() {
		try {
			const { data: status } = await TicketService.createStatusByBoardId(board_id, data);

			BoardStore.update((state) => {
				const boards = cloneDeep(state.boards);
				const board = boards.find((board) => board.id === status.board_id);

				if (!board) {
					return state;
				}

				board.statuses = [...board.statuses, status];
				state.boards = [...boards];

				return state;
			});

			DialogStore.close();
		} catch (error: any) {
			AlertStore.error(error);
		}
	}
</script>

<Dialog.Content>
	<Dialog.Header>
		<Dialog.Title>Add status</Dialog.Title>
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
		<Button type="submit" on:click={handleSubmit}>Add</Button>
	</Dialog.Footer>
</Dialog.Content>
