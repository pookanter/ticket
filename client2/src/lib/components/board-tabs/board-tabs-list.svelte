<script lang="ts">
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	type $$Props = TabsPrimitive.ListProps;

	let className: $$Props['class'] = undefined;
	export { className as class };
</script>

<div
	class="text-sm font-medium text-center text-gray-500 border-b border-gray-200 dark:text-gray-400 dark:border-gray-700"
>
	<TabsPrimitive.List class={cn('flex flex-wrap -mb-px', className)} {...$$restProps}>
		<slot />
	</TabsPrimitive.List>
</div>
