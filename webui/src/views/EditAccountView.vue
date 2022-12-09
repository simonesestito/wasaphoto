<script>
import PageSkeleton from "../components/PageSkeleton.vue";
import UsernameInput from "../components/UsernameInput.vue";
import {UsersService} from "../services";
import ErrorMsg from "../components/ErrorMsg.vue";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import {getCurrentUID} from "../services/auth-store";
import SuccessMsg from "../components/SuccessMsg.vue";

export default {
	name: "EditAccountView",
	components: {SuccessMsg, LoadingSpinner, ErrorMsg, UsernameInput, PageSkeleton},
	data() {
		return {
			loading: false,
			errorMessage: null,
			myProfile: null,
			success: false,
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
			this.success = false;
			this.loading = false;
			this.errorMessage = error;
		},
		async onUpdateUsername(username) {
			this.success = false;
			this.loading = true;
			this.errorMessage = null;
			try {
				await UsersService.setMyUsername(username);
				this.success = true;
			} catch (err) {
				this.errorMessage = err.toString();
			} finally {
				this.loading = false;
			}
		},
		async onUpdateDetails(event) {
			event.preventDefault();
			this.errorMessage = null;
			this.success = false;

			// Send HTTP!
			this.loading = true;
			try {
				await UsersService.setMyDetails({
					name: this.myProfile.name,
					surname: this.myProfile.surname,
					username: this.myProfile.username,
				});
				this.success = true;
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

		<SuccessMsg v-if="success" msg="Operation succeeded" />

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
						   :initial-input="myProfile.username" />
		</div>
	</PageSkeleton>
</template>
