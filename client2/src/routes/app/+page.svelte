<script lang="ts">
	import { AuthenService } from '$lib/services/authen-service';
	import { auth, me } from '$lib/store/authen';
	import { from } from 'rxjs';
	import { onMount } from 'svelte';
	import * as BoardTabs from '$lib/components/board-tabs/index';
	import { PlusOutline } from 'flowbite-svelte-icons';

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
			<BoardTabs.Trigger value="account">
				<span>Acccount</span>
			</BoardTabs.Trigger>
			<BoardTabs.Trigger value="password">
				<span>Password</span>
			</BoardTabs.Trigger>
		</BoardTabs.List>
		<BoardTabs.Content value="account">
			<h1>HELLO</h1>
		</BoardTabs.Content>
		<BoardTabs.Content value="password">
			<h1>WORLD</h1>
		</BoardTabs.Content>
	</BoardTabs.Root>
</section>
