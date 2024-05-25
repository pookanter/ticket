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
	class="text-sm font-medium text-center text-gray-500 border-b border-gray-200 dark:border-gray-700 dark:text-gray-400"
>
	<ul class="flex flex-wrap -mb-px">
		{#each boards as board}
			<Tab activate={board.id === target.id} {board} {onSelect} {onDelete} />
		{/each}
		<a
			href={null}
			class="inline-block p-4 border-b-2 border-transparent rounded-t-lg cursor-pointer hover:border-gray-300 hover:text-gray-600 dark:hover:text-gray-300"
			on:click={() => addBoard()}
		>
			<PlusOutline class="w-4 h-4" />
		</a>
	</ul>
</div>
