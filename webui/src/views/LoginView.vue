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
			keepSignedIn: true,
		};
	},
	methods: {
		async login(username) {
			this.loading = true;
			this.errorMessage = null;
			try {
				const {isNewUser} = await AuthService.doLogin(username, this.keepSignedIn);
				await router.push(isNewUser ? '/me/edit' : '/');
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
		<ErrorMsg v-if="this.errorMessage" :msg="this.errorMessage"/>
		<UsernameInput @submit="login" :loading="this.loading" submit-text="Login"/>

		<div class="form-check mt-3">
			<input class="form-check-input" type="checkbox" value="" id="keepSignedInCheck" v-model="keepSignedIn"
				   :checked="keepSignedIn">
			<label class="form-check-label" for="keepSignedInCheck">
				Keep me signed in
			</label>
		</div>
	</PageSkeleton>
</template>
