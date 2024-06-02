import { TicketService } from '$lib/services/ticket-service';
import { from } from 'rxjs';
import { writable, get } from 'svelte/store';

export type BoardStore = {
	initializing: boolean;
	boards: TicketService.Board[];
};

const BoardStore = writable<BoardStore>({
	initializing: true,
	boards: []
});

// async function createBoard(data: Parameters<typeof TicketService.createBoard>[0]) {
// 	from(TicketService.createBoard(data)).subscribe({
// 		next: ({ data: board }) => {
// 			BoardStore.update((store) => {
// 				store.boards.push(board);
// 				return store;
// 			});
// 		}
// 	});
// }

export { BoardStore };
