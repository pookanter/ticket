<script lang="ts">
	import StatusTab from './StatusTab.svelte';
	import type { Status } from './types';
	import * as Card from '$lib/components/ui/card';
	import { PlusIcon } from 'lucide-svelte';
	import Button from '$lib/components/ui/button/button.svelte';

	let statuses: Status[] = [
		{
			id: 1,
			name: 'Status 1'
		},
		{
			id: 2,
			name: 'Status 2'
		},
		{
			id: 3,
			name: 'Status 3'
		}
	];
	let target: Status = statuses[0];

	let onSelect = (status: Status) => {
		console.log('onSelect', status);
		target = status;
	};

	const dataDropDownId = 'dropdown';

	let onRename = () => {
		console.log('rename');
	};

	let onDelete = () => {
		statuses = statuses.filter((status) => status.id !== target.id);
	};

	const addStatus = () => {
		console.log('add status', dataDropDownId);
		const id = statuses.length + 1;
		statuses = [...statuses, { id, name: `Status ${id}` }];
	};
</script>

<div class="flex justify-start w-full gap-4">
	{#each statuses as status}
		<StatusTab data={status} {onSelect} {onDelete} />
	{/each}
	<Button
		variant="outline"
		class="box-content p-6"
		style="height: fit-content"
		on:click={() => addStatus()}><PlusIcon class="size-[18px] " /></Button
	>
</div>
