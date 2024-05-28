<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { fromEvent, type Subscription } from 'rxjs';
	import type { Status } from './types';
	import { onDestroy } from 'svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { EditOutline, TrashBinOutline, PenOutline } from 'flowbite-svelte-icons';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Input } from '$lib/components/ui/input/index.js';

	export let data: Status;
	export let onSelect: (data: Status) => void;
	export let onDelete: (data: Status) => void;

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
	function color() {
		
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

<DropdownMenu.Root>
	<DropdownMenu.Trigger disabled={isRenaming}>
		<Button
			variant="outline"
			class="justify-start w-64 p-4 text-start"
			style="height: fit-content"
			on:click={() => onSelect(data)}
		>
			{#if isRenaming}
				<Input
					bind:nativeElement={inputElement}
					bind:value={data.name}
					on:keydown={handleKeydown}
					class="py-0"
				/>
			{:else}
				<div class="flex items-center h-10">
					{data.name}
				</div>
			{/if}
		</Button>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content>
		<DropdownMenu.Item class="cursor-pointer" on:click={() => rename()}>
			<EditOutline class="w-4 h-4 mr-2" />
			<span>Rename</span>
		</DropdownMenu.Item>
		<DropdownMenu.Sub on:click={() => rename()}>
			<DropdownMenu.SubTrigger>
				<PenOutline class="w-4 h-4 mr-2" />
				<span>Color</span>
			</DropdownMenu.SubTrigger>
			<DropdownMenu.SubContent class="w-[175px]">
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(250, 175, 168);"></div>
					<span>Coral Red Orange</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(243, 159, 118);"></div>
					<span>Orange Peach</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(255, 248, 184);"></div>
					<span>Pink Sand</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(226, 246, 211);"></div>
					<span>Green Mint</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(180, 221, 211);"></div>
					<span>Green Sage</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(212, 228, 237);"></div>
					<span>Gray Fox</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(174, 204, 220);"></div>
					<span>Gray Storm</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(211, 191, 219);"></div>
					<span>Purple Dusk</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(246, 226, 221);"></div>
					<span>Pink Blossom</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(233, 227, 212);"></div>
					<span>Brown Clay</span>
				</DropdownMenu.Item>
				<DropdownMenu.Item class="cursor-pointer">
					<div class="w-4 h-4 mr-2" style="background-color: rgb(239, 239, 241);"></div>
					<span>Gray Shalk</span>
				</DropdownMenu.Item>
			</DropdownMenu.SubContent>
		</DropdownMenu.Sub>
		<DropdownMenu.Item class="cursor-pointer" on:click={() => onDelete(data)}>
			<TrashBinOutline class="w-4 h-4 mr-2" />
			<span>Delete</span>
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
