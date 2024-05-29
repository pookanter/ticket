<script>
	import { AuthenService } from '$lib/services/authen-service';
	import { auth, me } from '$lib/store/authen';
	import { from } from 'rxjs';
	import { onMount } from 'svelte';

	onMount(() => {
		from(AuthenService.getMe()).subscribe({
			next: ({ data }) => {
				me.set(data);
			},
			error: (error) => {
				AuthenService.removeAuthorization();
				me.set(null);
				auth.set(null);
			}
		});
	});
</script>
