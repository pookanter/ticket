import { TicketService } from '$lib/services/ticket-service';
import { from } from 'rxjs';
import { writable, get } from 'svelte/store';

export type TicketStore = {
	initializing: boolean;
	boards: TicketService.Board[];
};

const ticketStore = writable<TicketStore>({
	initializing: true,
	boards: []
});

function createBoard(data: Parameters<typeof TicketService.createBoard>[0]) {
	from(TicketService.createBoard(data)).subscribe({
		next: ({ data: board }) => {
			ticketStore.update((store) => {
				store.boards.push(board);
				return store;
			});
		}
	});
}

export { ticketStore, createBoard };
