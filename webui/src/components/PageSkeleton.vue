<template>
	<div class="mb-4">
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">{{ title }}</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button v-for="action in actions" :key="action.text" type="button"
							class="btn btn-sm btn-outline-secondary" @click="action.onClick">
						{{ action.text }}
					</button>
					<button v-if="mainAction" type="button"
							class="btn btn-sm btn-outline-primary" @click="mainAction.onClick">
						{{ mainAction.text }}
					</button>
				</div>
			</div>
		</div>

		<slot />
	</div>
</template>

<script>
export default {
	name: "PageSkeleton",
	props: [
		'title',
		'actions',
		'mainAction',
	],
};
</script>
