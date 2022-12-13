<script>
import {AuthService} from "../services";
import UsernameInput from '../components/UsernameInput.vue';
import router from "../router";
import PageSkeleton from "../components/PageSkeleton.vue";
import {getCurrentUID} from "../services/auth-store";

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
				const previousPath = this.$route.query.previous || '/';
				await router.push(isNewUser ? '/me/edit' : previousPath);
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
	},
	mounted() {
		if (getCurrentUID() != null) {
			// Already logged in, redirect to My Profile
			router.replace(this.$route.query.previous || '/me');
		} else if (this.$route.query.previous) {
			this.errorMessage = 'Login to continue';
		}
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
