<script>
import {AuthService} from "../services";
import UsernameInput from '../components/UsernameInput.vue';
import router from "../router";

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
		UsernameInput,
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Login</h1>
		</div>

		<UsernameInput @submit="login" @error="onError" :loading="this.loading" />

		<ErrorMsg v-if="this.errorMessage" :msg="this.errorMessage" />
	</div>
</template>
