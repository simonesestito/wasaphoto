<script>
import UsernameInput from "../components/UsernameInput.vue";
import {UsersService} from "../services";
import UsersList from '../components/UsersList.vue';
import PageSkeleton from "../components/PageSkeleton.vue";

export default {
	components: {UsersList, UsernameInput, PageSkeleton},
	data: function () {
        return {
			searchedUsername: null,
			errorMessage: null,
        };
    },
    methods: {
        async refresh(username) {
			this.searchedUsername = username;
        },
		async loadNextPage(cursor) {
			return UsersService.searchUsers(this.searchedUsername, cursor);
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
		<UsernameInput submit-text="Search" @submit="refresh" />

		<!-- Error -->
		<ErrorMsg v-if="errorMessage" :msg="errorMessage" />

		<!-- Show users -->
		<UsersList :loader-function="loadNextPage" :refresh-key="this.searchedUsername" />
	</PageSkeleton>
</template>
