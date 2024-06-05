import { TicketService } from '$lib/services/ticket-service';
import { cloneDeep } from 'lodash';
import { from } from 'rxjs';
import { writable, get } from 'svelte/store';

export type BoardState = {
	selected?: TicketService.BoardFullDetail;
	boards: TicketService.Board[];
};

function defaultState(): BoardState {
	return {
		selected: undefined,
		boards: []
	};
}

const boardStore = writable<BoardState>(defaultState());

function selectBoard(board: TicketService.BoardFullDetail) {
	boardStore.update((state) => {
		return { ...state, selected: board };
	});
}

function addBoard({ board }: { board: TicketService.Board }) {
	boardStore.update((state) => {
		const boards = cloneDeep(state.boards);
		boards.push(board);

		state.boards = boards;

		return state;
	});
}

function updateBoard({ board }: { board: TicketService.Board }) {
	boardStore.update((state) => {
		const boards = cloneDeep(state.boards);

		const index = boards.findIndex((b) => b.id === board.id);
		if (index !== -1) {
			boards[index] = board;
		}

		state.boards = boards;

		return state;
	});
}

function addStatus({ status }: { status: TicketService.Status }) {
	boardStore.update((state) => {
		const selected = cloneDeep(state.selected);
		if (selected && selected.id === status.board_id) {
			selected.statuses.push(status);
		}

		state.selected = selected;

		return state;
	});
}

function updateStatus({ status }: { status: TicketService.Status }) {
	boardStore.update((state) => {
		const selected = cloneDeep(state.selected);
		const index = selected.statuses.findIndex((s) => s.id === status.id);
		if (index !== -1) {
			selected.statuses[index] = status;
		}

		state.selected = selected;

		return state;
	});
}

function addTicket({ board_id, ticket }: { board_id: number; ticket: TicketService.Ticket }) {
	boardStore.update((state) => {
		const selected = cloneDeep(state.selected);
		if (selected && selected.id === board_id) {
			const status = selected.statuses.find((s) => s.id === ticket.status_id);
			if (status) {
				status.tickets.push(ticket);
			}
		}

		state.selected = selected;

		return state;
	});
}

function updateTicket({ ticket }: { ticket: TicketService.Ticket }) {
	boardStore.update((state) => {
		const selected = cloneDeep(state.selected);
		const status = selected.statuses.find((s) => s.id === ticket.status_id);
		if (status) {
			const index = status.tickets.findIndex((t) => t.id === ticket.id);
			if (index !== -1) {
				status.tickets[index] = ticket;
			}
		}

		state.selected = selected;

		return state;
	});
}

export const BoardStore = {
	...boardStore,
	defaultState,
	selectBoard,
	addBoard,
	updateBoard,
	addStatus,
	updateStatus,
	addTicket,
	updateTicket
};
