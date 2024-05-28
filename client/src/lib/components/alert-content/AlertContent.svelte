<script lang="ts">
	import * as AlertDialog from '$lib/components/ui/alert-dialog';

	export let title: string;
	export let description: string;
</script>

<AlertDialog.Content>
	<AlertDialog.Header>
		<AlertDialog.Title>{title}</AlertDialog.Title>
		<AlertDialog.Description>
			{description}
		</AlertDialog.Description>
	</AlertDialog.Header>
	<AlertDialog.Footer>
		<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
		<AlertDialog.Action>Ok</AlertDialog.Action>
	</AlertDialog.Footer>
</AlertDialog.Content>
