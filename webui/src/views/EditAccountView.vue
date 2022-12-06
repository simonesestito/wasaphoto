<script>
import PageSkeleton from "../components/PageSkeleton.vue";
import UsernameInput from "../components/UsernameInput.vue";
import {UsersService} from "../services";
import ErrorMsg from "../components/ErrorMsg.vue";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import {getCurrentUID} from "../services/auth-store";

export default {
	name: "EditAccountView",
	components: {LoadingSpinner, ErrorMsg, UsernameInput, PageSkeleton},
	data() {
		return {
			loading: false,
			errorMessage: null,
			myProfile: null,
		};
	},
	methods: {
		async loadUser() {
			this.loading = true;
			try {
				this.myProfile = await UsersService.getUserProfile(getCurrentUID());
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
		onError(error) {
			this.loading = false;
			this.errorMessage = error;
		},
		async onUpdateUsername(username) {
			this.loading = true;
			this.errorMessage = null;
			try {
				await UsersService.setMyUsername(username);
				// TODO: Show OK alert
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
		async onUpdateDetails(event) {
			event.preventDefault();
			this.errorMessage = null;

			// Send HTTP!
			this.loading = true;
			try {
				await UsersService.setMyDetails({
					name: this.myProfile.name,
					surname: this.myProfile.surname,
					username: this.myProfile.username,
				});
				// TODO: OK message
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		}
	},
	mounted() {
		this.loadUser();
	}
}
</script>

<template>
	<PageSkeleton title="Edit account">

		<div v-if="errorMessage">
			<ErrorMsg :msg="errorMessage"/>
			<hr>
		</div>

		<LoadingSpinner v-if="loading"/>

		<div v-if="myProfile">
			<h3>Edit account data</h3>
			<form @submit="onUpdateDetails">
				<div class="mb-3">
					<label for="nameInput" class="form-label">Name</label>
					<input autocomplete="given-name" class="form-control" id="nameInput" v-model="myProfile.name"
						   minlength="3"
						   maxlength="256" required>
				</div>
				<div class="mb-3">
					<label for="surnameInput" class="form-label">Surname</label>
					<input autocomplete="family-name" class="form-control" id="surnameInput" v-model="myProfile.surname"
						   maxlength="256">
				</div>

				<button type="submit" class="btn btn-primary mb-3" :disabled="loading">Update account data</button>
			</form>

			<hr>

			<h3>Edit username</h3>
			<UsernameInput submit-text="Change username" :loading="this.loading" @submit="onUpdateUsername"
						   @error="onError" :initial-input="myProfile.username" />
		</div>
	</PageSkeleton>
</template>
