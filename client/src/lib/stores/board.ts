import { TicketService } from '$lib/services/ticket-service';
import { cloneDeep } from 'lodash';
import { from } from 'rxjs';
import { writable, get } from 'svelte/store';

export type BoardState = {
	initializing: boolean;
	selected: TicketService.Board | null;
	boards: TicketService.Board[];
};

function defaultState(): BoardState {
	return {
		initializing: true,
		selected: null,
		boards: []
	};
}

const boardStore = writable<BoardState>(defaultState());

function addBoard({ board }: { board: TicketService.Board }) {
	boardStore.update((store) => {
		const boards = cloneDeep(store.boards);
		boards.push(board);

		store.boards = boards;

		return store;
	});
}

function addStatus({ status }: { status: TicketService.Status }) {
	boardStore.update((store) => {
		const boards = cloneDeep(store.boards);
		const board = boards.find((b) => b.id === status.board_id);
		if (board) {
			board.statuses.push(status);

			if (store.selected && store.selected.id === board.id) {
				store.selected = board;
			}
		}

		store.boards = boards;

		return store;
	});
}

function addTicket({ board_id, ticket }: { board_id: number; ticket: TicketService.Ticket }) {
	boardStore.update((store) => {
		const boards = cloneDeep(store.boards);
		const board = boards.find((b) => b.id === board_id);
		if (board) {
			const status = board.statuses.find((s) => s.id === ticket.status_id);
			if (status) {
				status.tickets.push(ticket);
			}

			if (store.selected && store.selected.id === board.id) {
				store.selected = board;
			}
		}

		store.boards = boards;

		return store;
	});
}

export const BoardStore = { ...boardStore, defaultState, addBoard, addStatus, addTicket };
