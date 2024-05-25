<script lang="ts">
	import Tab from './Tab.svelte';
	import type { Board } from './types';
	import { PlusOutline } from 'flowbite-svelte-icons';

	let boards: Board[] = [
		{
			id: 1,
			name: 'Board 1'
		},
		{
			id: 2,
			name: 'Board 2'
		},
		{
			id: 3,
			name: 'Board 3'
		}
	];
	let target: Board = boards[0];

	let onSelect = (board: Board) => {
		target = board;
	};

	const dataDropDownId = 'dropdown';

	let onRename = () => {
		console.log('rename');
	};

	let onDelete = () => {
		boards = boards.filter((board) => board.id !== target.id);
	};

	const addBoard = () => {
		console.log('add board', dataDropDownId);
		const id = boards.length + 1;
		boards = [...boards, { id, name: `Board ${id}` }];
	};
</script>

<div
	class="border-b border-gray-200 text-center text-sm font-medium text-gray-500 dark:border-gray-700 dark:text-gray-400"
>
	<ul class="-mb-px flex flex-wrap">
		{#each boards as board}
			<Tab activate={board.id === target.id} {board} {onSelect} {onDelete} />
		{/each}
		<a
			href={null}
			class="inline-block cursor-pointer rounded-t-lg border-b-2 border-transparent p-4 hover:border-gray-300 hover:text-gray-600 dark:hover:text-gray-300"
			on:click={() => addBoard()}
		>
			<PlusOutline class="h-4 w-4" />
		</a>
	</ul>
</div>
