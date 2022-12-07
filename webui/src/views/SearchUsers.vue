<script>
import PageSkeleton from "../components/PageSkeleton.vue";
import UsernameInput from "../components/UsernameInput.vue";
import {UsersService} from "../services";
import ErrorMsg from "../components/ErrorMsg.vue";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import UserListItem from "../components/UserListItem.vue";
import ShowMore from '../components/ShowMore.vue';

export default {
	components: {UserListItem, LoadingSpinner, ErrorMsg, UsernameInput, PageSkeleton, ShowMore},
	data: function () {
        return {
            errorMessage: null,
            loading: false,
            foundUsers: [],
			pageCursor: null,
			searchedUsername: null,
        };
    },
    methods: {
        async refresh(username) {
			this.searchedUsername = username;
			this.pageCursor = null;
			this.foundUsers = [];
            await this.loadNextPage();
        },
		async loadNextPage() {
			this.loading = true;
			this.errorMessage = null;

			try {
				const response = await UsersService.searchUsers(this.searchedUsername, this.pageCursor);
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
    }
}
</script>

<template>
	<PageSkeleton title="Search Users">
		<!-- Username search bar -->
		<UsernameInput submit-text="Search" @error="onError" @submit="refresh" />

		<!-- Error -->
		<ErrorMsg v-if="errorMessage" :msg="errorMessage" />

		<!-- Show users -->
		<UserListItem v-for="user in foundUsers" @error="onError" :user="user" />

		<!-- Loading -->
		<LoadingSpinner v-if="loading" />

		<!-- Show more -->
		<ShowMore v-if="pageCursor" @loadNext="loadNextPage" />
	</PageSkeleton>
</template>
