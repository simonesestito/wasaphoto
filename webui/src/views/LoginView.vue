<script>
import {AuthService} from "../services";
import UsernameInput from '../components/UsernameInput.vue';
import router from "../router";
import PageSkeleton from "../components/PageSkeleton.vue";

export default {
	data: function () {
		return {
			errorMessage: null,
			loading: false,
		};
	},
	methods: {
		async login(username) {
			this.loading = true;
			this.errorMessage = null;
			try {
				const {isNewUser} = await AuthService.doLogin(username);
				await router.push(isNewUser ? '/me' : '/');
			} catch (e) {
				this.errorMessage = e.toString();
			} finally {
				this.loading = false;
			}
		},
		onError(error) {
			this.errorMessage = error;
		}
	},
	components: {
		PageSkeleton,
		UsernameInput,
	}
}
</script>

<template>
	<PageSkeleton title="Login">
		<UsernameInput @submit="login" @error="onError" :loading="this.loading" submit-text="Login" />
		<ErrorMsg v-if="this.errorMessage" :msg="this.errorMessage" />
	</PageSkeleton>
</template>
