<script>
export default {
	emits: [
		'submit',
		'error',
	],
	props: [
		'loading',
	],
	data() {
		return {
			userInput: '',
			errorMessage: null,
		};
	},
	methods: {
		login(event) {
			event.preventDefault();
			if (this.validate())
				this.$emit('submit', this.userInput);
			else
				this.$emit('error', this.errorMessage);
		},
		validate() {
			if (this.userInput.length < 3)
				this.errorMessage = 'Username too short';
			else if (this.userInput.length > 16)
				this.errorMessage = 'Username too long';
			else if (!this.userInput.match(/^[a-z_0-9]+$/))
				this.errorMessage = 'Username not valid';
			else
				this.errorMessage = null;

			return this.errorMessage == null;
		},
	},
};
</script>

<template>
	<form @submit="login">
		<div class="mb-3">
			<label for="usernameInput" class="form-label">Username</label>
			<input type="text" class="form-control" id="usernameInput" v-model="userInput">
		</div>

		<input class="btn btn-primary" type="submit" value="Login" :disabled="loading">
	</form>

</template>
