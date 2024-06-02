import { TicketService } from '$lib/services/ticket-service';
import { from } from 'rxjs';
import { writable, get } from 'svelte/store';

export type BoardState = {
	initializing: boolean;
	selected?: TicketService.Board;
	boards: TicketService.Board[];
};

function defaultState(): BoardState {
	return {
		initializing: true,
		selected: undefined,
		boards: []
	};
}

const boardStore = writable<BoardState>(defaultState());

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

export const BoardStore = { ...boardStore, defaultState };
