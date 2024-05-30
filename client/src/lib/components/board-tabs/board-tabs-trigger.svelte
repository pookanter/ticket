<script lang="ts">
	import {
		DotsHorizontalOutline,
		DotsVerticalOutline,
		EditOutline,
		TrashBinOutline
	} from 'flowbite-svelte-icons';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { Tabs as TabsPrimitive, type CustomEventHandler } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import { clickOutsideAction } from '$lib/directives/click-outside';
	import Button from '$lib/components/ui/button/button.svelte';

	type $$Props = TabsPrimitive.TriggerProps;
	type $$Events = TabsPrimitive.TriggerEvents;

	let className: $$Props['class'] = undefined;
	export let value: $$Props['value'];
	export { className as class };

	let menuEnabled = false;
	function onFocus(e: CustomEventHandler<FocusEvent, HTMLButtonElement>) {
		setTimeout(() => {
			menuEnabled = true;
		}, 0);
	}
</script>

<TabsPrimitive.Trigger
	class={cn(
		'inline-block py-2 px-4 border-b-2 border-transparent rounded-t-lg !cursor-pointer hover:text-gray-600 hover:border-gray-300 dark:hover:text-gray-300 data-[state=active]:text-blue-600  data-[state=active]:border-blue-600',
		className
	)}
	{value}
	{...$$restProps}
	on:click
	on:focus={onFocus}
>
	<DropdownMenu.Root>
		<div class="flex justify-between">
			<div>
				<slot />
			</div>
			<DropdownMenu.Trigger>
				<div class="p-1 ml-4 rounded cursor-pointer hover:text-accent-foreground hover:bg-accent">
					<DotsHorizontalOutline class="size-4" />
				</div>
			</DropdownMenu.Trigger>
		</div>
		<DropdownMenu.Content>
			<DropdownMenu.Item class="cursor-pointer">
				<EditOutline class="mr-2 size-6" />
				<span>Rename</span>
			</DropdownMenu.Item>
			<DropdownMenu.Item class="cursor-pointer">
				<TrashBinOutline class="w-4 h-4 mr-2" />
				<span>Delete</span>
			</DropdownMenu.Item>
		</DropdownMenu.Content>
	</DropdownMenu.Root>
</TabsPrimitive.Trigger>
