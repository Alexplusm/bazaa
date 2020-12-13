import Vue from 'vue';
import Vuex from 'vuex';
import { api } from '../api';
import { login } from '../domain/storage-domain'

Vue.use(Vuex);

export const store = new Vuex.Store({
	state: {
		auth: {
			authorized: false,
		},
		extSystems: [],
		currentExtSystem: null,
		currentGame: null,
		games: [],
	},
	mutations: {
		authorize(state) {
			login();
			state.auth.authorized = true;
		},
		setExtSystemList(state, extSystems) {
			state.extSystems = extSystems;
		},
		setCurrentExtSystem(state, extSystem) {
			state.currentExtSystem = extSystem;
		},
		setGames(state, games) {
			state.games = games;
		},
		setCurrentGame(state, game) {
			state.currentGame = game;
		},
	},
	actions: {
		authorize({commit}) {
			commit('authorize');
		},
		getExtSystemList({ commit }) {
			return api.extSystem
				.list()
				.then((data) => commit('setExtSystemList', data.extSystems));
		},
		createExtSystem(context, extSystem) {
			return api.extSystem.create(extSystem);
		},
		setCurrentExtSystem({ commit, dispatch }, extSystem) {
			commit('setCurrentExtSystem', extSystem);
			dispatch('getGameList');
		},
		getGameList({ commit, state }) {
			const { currentExtSystem } = state;

			if (currentExtSystem && currentExtSystem.extSystemId) {
				return api.game
					.list(state.currentExtSystem.extSystemId)
					.then((data) => commit('setGames', data.games));
			}
		},
		getGameDetails({ commit, state }, gameId) {
			const { currentExtSystem } = state;

			if (currentExtSystem && currentExtSystem.extSystemId) {
				return api.game
					.details(gameId, currentExtSystem.extSystemId)
					.then((data) => commit('setCurrentGame', data));
			}
		},
		updateGameWithArchive(
			{ dispatch },
			{ gameId, file, progressCallback }
		) {
			return api.game
				.updateWithFile(gameId, file, progressCallback)
				.finally(() => dispatch('getGameDetails', gameId));
		},
		createGame(_, game) {
			return api.game.create(game);
		},
	},
	getters: {
		extSystems(state) {
			return state.extSystems;
		},
		currentExtSystem(state) {
			return state.currentExtSystem;
		},
		games(state) {
			return state.games;
		},
		currentGame(state) {
			return state.currentGame;
		},
	},
	modules: {},
});
