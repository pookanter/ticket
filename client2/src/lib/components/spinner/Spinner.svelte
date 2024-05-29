<script lang="ts">
	import { cn } from '$lib/utils';

	let className = '';

	export { className as class };
</script>

<div
	class={cn(
		'text-surface inline-block h-4 w-4 animate-spin rounded-full border-2 border-solid border-current border-e-transparent align-[-0.125em] motion-reduce:animate-[spin_1.5s_linear_infinite] dark:text-white',
		className
	)}
	role="status"
>
	<span
		class="!absolute !-m-px !h-px !w-px !overflow-hidden !whitespace-nowrap !border-0 !p-0 ![clip:rect(0,0,0,0)]"
		>Loading...</span
	>
</div>
