<script>
import {FollowService, UsersService} from "../services";
import UsersList from '../components/UsersList.vue';
import PageSkeleton from "../components/PageSkeleton.vue";

export default {
	components: {UsersList, PageSkeleton},
	data: function () {
		return {
			errorMessage: null,
			username: null,
		};
	},
	methods: {
		refresh(username) {
			this.user = null;
			this.username = username;
		},
		// Used in <UsersList> component as the refresh function. It returns a value.
		async loadNextPage(cursor) {
			if (this.user == null) {
				this.user = await UsersService.getByUsername(this.username);
			}
			return FollowService.listFollowers(this.user.id, cursor);
		},
		onError(err) {
			this.errorMessage = err.toString();
		},
	},
	computed: {
		pageTitle() {
			return `@${this.username}'s followers`
		},
	},
	mounted() {
		this.refresh(this.$route.params.username);
	},
	created() {
		this.$watch(
			() => this.$route.params,
			(params, _) => {
				// react to route changes...
				this.refresh(params.username);
			}
		)
	},
}
</script>

<template>
	<PageSkeleton :title="pageTitle">
		<!-- Error -->
		<ErrorMsg v-if="errorMessage" :msg="errorMessage" />

		<!-- Show users -->
		<UsersList :loader-function="loadNextPage" :refresh-key="username" />
	</PageSkeleton>
</template>
