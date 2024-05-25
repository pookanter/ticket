<script lang="ts">
	import type { Board } from './types';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { EditOutline, TrashBinOutline } from 'flowbite-svelte-icons';
	import { Input } from '$lib/components/ui/input/index.js';
	import { fromEvent, type Subscription } from 'rxjs';
	import { onDestroy } from 'svelte';

	export let activate: boolean;
	export let data: Board;
	export let onSelect: (data: Board) => void;
	export let onDelete: (data: Board) => void;

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

	onDestroy(() => {
		if (click$) {
			click$.unsubscribe();
		}
	});
</script>

<li class="me-2">
	{#if activate}
		<DropdownMenu.Root>
			<DropdownMenu.Trigger disabled={isRenaming}>
				{#if isRenaming}
					<Input
						bind:nativeElement={inputElement}
						bind:value={data.name}
						on:keydown={handleKeydown}
						style="width: {data.name.length + 3}ch"
						class="mt-[6px]"
					/>
				{:else}
					<div
						class="relative inline-block p-4 border-b-2 rounded-t-lg cursor-pointer active border-primary-600 text-primary-600 dark:border-primary-500 dark:text-blprimaryue-500"
						aria-current="page"
					>
						{data.name}
					</div>
				{/if}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content>
				<DropdownMenu.Item class="cursor-pointer" on:click={() => rename()}>
					<EditOutline class="w-4 h-4 mr-2" />
					<span>Rename</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer" on:click={() => onDelete(data)}>
					<TrashBinOutline class="w-4 h-4 mr-2" />
					<span>Delete</span>
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	{:else}
		<a
			href={null}
			class="inline-block p-4 border-b-2 border-transparent rounded-t-lg cursor-pointer hover:border-gray-300 hover:text-gray-600 dark:hover:text-gray-300"
			on:click={() => onSelect(data)}>{data.name}</a
		>
	{/if}
</li>
