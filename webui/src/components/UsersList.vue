<script>
import ErrorMsg from "../components/ErrorMsg.vue";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import UserListItem from "../components/UserListItem.vue";
import ShowMore from '../components/ShowMore.vue';

export default {
	name: 'UsersList',
	components: {UserListItem, LoadingSpinner, ErrorMsg, ShowMore},
	props: ['refreshKey', 'loaderFunction'],
	data() {
		return {
			errorMessage: null,
			loading: false,
			foundUsers: [],
			pageCursor: null,
		};
	},
	methods: {
		async refresh() {
			this.pageCursor = null;
			this.foundUsers = [];
			await this.loadNextPage();
		},
		async loadNextPage() {
			// Handle double loading or null (/undefined) parameters.
			if (!this.refreshKey || this.loading) {
				return;
			}

			this.loading = true;
			this.errorMessage = null;

			try {
				const response = await this.loaderFunction(this.pageCursor);
				this.foundUsers.push(...response.pageData);
				this.pageCursor = response.nextPageCursor;
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
		onError(err) {
			this.errorMessage = err.toString();
		},
	},
	watch: {
		refreshKey() {
			this.refresh();
		},
	},
	mounted() {
		this.refresh();
	},
}
</script>

<template>
	<!-- Error -->
	<ErrorMsg v-if="errorMessage" :msg="errorMessage"/>

	<!-- Show users -->
	<UserListItem v-for="user in foundUsers" @error="onError" :user="user" :key="user.id"/>

	<!-- Loading -->
	<LoadingSpinner v-if="loading"/>

	<!-- Show more -->
	<ShowMore v-if="pageCursor && !loading" @loadMore="loadNextPage"/>
</template>
