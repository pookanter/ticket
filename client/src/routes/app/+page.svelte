<script lang="ts">
	import BoardTabset from '$lib/components/board-tabset/BoardTabset.svelte';
	import * as Card from '$lib/components/ui/card';
	// import { Card } from '$lib/components/ui/card/index.js';
</script>

<section class="size-full">
	<div class="flex flex-col max-w-screen-xl px-4 py-8 mx-auto left-o size-full">
		<BoardTabset />
		<div class="overflow-x-scroll size-full">
			<Card.Root class="w-[350px]">
				<Card.Header>
					<Card.Title>Create project</Card.Title>
					<Card.Description>Deploy your new project in one-click.</Card.Description>
				</Card.Header>
				<Card.Content>
					<form>
						<div class="grid items-center w-full gap-4">
							<div class="flex flex-col space-y-1.5"></div>
							<div class="flex flex-col space-y-1.5"></div>
						</div>
					</form>
				</Card.Content>
			</Card.Root>
		</div>
	</div>
</section>
