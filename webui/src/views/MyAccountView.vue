<script>
import PageSkeleton from "../components/PageSkeleton.vue";
import {getCurrentUID, saveAuthToken} from "../services/auth-store";
import router from "../router";
import {UsersService} from "../services";
import LoadingSpinner from "../components/LoadingSpinner.vue";

export default {
	components: {LoadingSpinner, PageSkeleton},
	data: function () {
		return {
			errorMessage: null,
			loading: false,
			myProfile: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errorMessage = null;
			try {
				this.myProfile = await UsersService.getUserProfile(getCurrentUID());
			} catch (e) {
				this.errorMessage = e.toString();
			} finally {
				this.loading = false;
			}
		},
		async onClick() {
			saveAuthToken(null);
			await router.replace('/login');
		},
		async edit() {
			await router.push('/me/edit');
		},
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<PageSkeleton title="My Account" :main-action="{text:'Logout', onClick: this.onClick}"
				  :actions="[{text:'Edit account', onClick: this.edit}]">
		<ErrorMsg v-if="errorMessage" :msg="errorMessage"/>

		<LoadingSpinner v-if="loading"/>

		<!-- User simple summary -->
		<div v-if="myProfile">
			<p>Hi, {{ myProfile.name }} {{ myProfile.surname }} (@{{ myProfile.username }})</p>
			<p>
				<RouterLink :to="'/users/' + myProfile.username">Show my profile</RouterLink>
			</p>
		</div>
	</PageSkeleton>
</template>

<style>

</style>
