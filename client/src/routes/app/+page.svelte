<script lang="ts">
	import BoardTabset from '$lib/components/board-tabset/BoardTabset.svelte';
	import StausTabset from '$lib/components/status-tabset/StausTabset.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import ScrollArea from '$lib/components/ui/scroll-area/scroll-area.svelte';
	import { PlusIcon } from 'lucide-svelte';

	type Status = {
		name: string;
	};
	const stauses: Status[] = [{ name: 'To Do' }, { name: 'In Progress' }, { name: 'Done' }];
</script>

<section class="h-[90vh] w-full">
	<div class="flex flex-col max-w-screen-xl px-4 py-8 mx-auto left-o size-full">
		<BoardTabset />
		<ScrollArea class="my-4 size-full" orientation="horizontal">
			<StausTabset />
		</ScrollArea>
	</div>
</section>
