<script>
	import SignIn from './SignIn.svelte';
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section class="m-auto flex h-[90vh] items-center justify-center">
	<div class="flex h-full flex-shrink-[0.6] flex-col items-center justify-center">
		<div class="grid max-w-screen-xl px-4 py-8 mx-auto lg:grid-cols-12 lg:gap-8 lg:py-16 xl:gap-0">
			<div class="mr-auto place-self-center lg:col-span-7">
				<h1
					class="max-w-2xl mb-4 text-4xl font-extrabold leading-none tracking-tight md:text-5xl xl:text-6xl dark:text-white"
				>
					Streamline Your Support Tickets with Ease
				</h1>
				<p
					class="max-w-2xl mb-6 font-light text-gray-500 md:text-lg lg:mb-8 lg:text-xl dark:text-gray-400"
				>
					Transform the way you manage customer support. Our Kanban-style ticketing system
					simplifies workflows, boosts productivity, and ensures nothing slips through the cracks.
				</p>
			</div>
			<div class="flex justify-center mt-8 lg:col-span-5 lg:mt-0 lg:justify-end">
				<SignIn />
			</div>
		</div>
	</div>
</section>
