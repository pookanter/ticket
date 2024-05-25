<script lang="ts">
	import type { Board } from './types';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { EditOutline, TrashBinOutline } from 'flowbite-svelte-icons';
	import { Input } from '$lib/components/ui/input/index.js';
	import { fromEvent, type Subscription } from 'rxjs';

	export let activate: boolean;
	export let board: Board;
	export let onSelect: (board: Board) => void;
	export let onDelete: (board: Board) => void;

	let click$: Subscription;

	let inputElement: any;
	let isRenaming = false;
	function rename() {
		if (click$) {
			click$.unsubscribe();
		}

		click$ = fromEvent(document, 'mousedown', {
			capture: true
		}).subscribe((ev) => {
			const ele = ev.target as HTMLElement;

			if (inputElement && !inputElement.contains(ele) && !ev.defaultPrevented) {
				isRenaming = false;
				click$.unsubscribe();
			}
		});

		isRenaming = true;
		setTimeout(() => {
			if (inputElement && inputElement.focus) {
				inputElement.focus();
			}
		}, 100);
	}
	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			e.preventDefault();
			e.stopPropagation();
			isRenaming = false;
		}
	}
</script>

<li class="me-2">
	{#if activate}
		<DropdownMenu.Root>
			<DropdownMenu.Trigger disabled={isRenaming}>
				{#if isRenaming}
					<Input
						bind:nativeElement={inputElement}
						bind:value={board.name}
						on:keydown={handleKeydown}
						style="width: {board.name.length + 3}ch"
						class="mt-[6px]"
					/>
				{:else}
					<div
						class="active border-primary-600 text-primary-600 dark:border-primary-500 dark:text-blprimaryue-500 relative inline-block cursor-pointer rounded-t-lg border-b-2 p-4"
						aria-current="page"
					>
						{board.name}
					</div>
				{/if}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content>
				<DropdownMenu.Item class="cursor-pointer" on:click={() => rename()}>
					<EditOutline class="mr-2 h-4 w-4" />
					<span>Rename</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer" on:click={() => onDelete(board)}>
					<TrashBinOutline class="mr-2 h-4 w-4" />
					<span>Delete</span>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	{:else}
		<a
			href={null}
			class="inline-block cursor-pointer rounded-t-lg border-b-2 border-transparent p-4 hover:border-gray-300 hover:text-gray-600 dark:hover:text-gray-300"
			on:click={() => onSelect(board)}>{board.name}</a
		>
	{/if}
</li>
