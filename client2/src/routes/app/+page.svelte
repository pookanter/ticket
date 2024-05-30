<script lang="ts">
	import { AuthenService } from '$lib/services/authen-service';
	import { auth, me } from '$lib/store/authen';
	import { from } from 'rxjs';
	import { onMount } from 'svelte';
	import * as BoardTabs from '$lib/components/board-tabs/index';

	onMount(() => {
		// from(AuthenService.getMe()).subscribe({
		// 	next: ({ data }) => {
		// 		me.set(data);
		// 	},
		// 	error: (error) => {
		// 		AuthenService.removeAuthorization();
		// 		me.set(null);
		// 		auth.set(null);
		// 	}
		// });
	});
</script>

<section class="h-[90vh] w-full">
	<BoardTabs.Root>
		<BoardTabs.List>
			<BoardTabs.Trigger>
				<span>Board 1</span>
			</BoardTabs.Trigger>
			<BoardTabs.Trigger>
				<span>Board 2</span>
			</BoardTabs.Trigger>
		</BoardTabs.List>
		<!-- <BoardTabs.Content>
			<BoardTabs.Tab>
				<BoardTabs.TabTrigger>
					<BoardTabs.TabIcon />
					<BoardTabs.TabTitle>Board 1</BoardTabs.TabTitle>
				</BoardTabs.TabTrigger>
				<BoardTabs.TabContent>
					<BoardTabs.TabContentTrigger>
						<BoardTabs.TabContentIcon />
						<BoardTabs.TabContentTitle>Board 1 Content</BoardTabs.TabContentTitle>
					</BoardTabs.TabContentTrigger>
				</BoardTabs.TabContent>
			</BoardTabs.Tab>
			<BoardTabs.Tab>
				<BoardTabs.TabTrigger>
					<BoardTabs.TabIcon />
					<BoardTabs.TabTitle>Board 2</BoardTabs.TabTitle>
				</BoardTabs.TabTrigger>
				<BoardTabs.TabContent>
					<BoardTabs.TabContentTrigger>
						<BoardTabs.TabContentIcon />
						<BoardTabs.TabContentTitle>Board 2 Content</BoardTabs.TabContentTitle>
					</BoardTabs.TabContentTrigger>
				</BoardTabs.TabContent>
			</BoardTabs.Tab>
		</BoardTabs.Content> -->
	</BoardTabs.Root>
</section>
