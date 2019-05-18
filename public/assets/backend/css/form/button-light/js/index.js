import FluentRevealEffect from "../lib/js/main";

// console.log(FluentRevealEffect)

FluentRevealEffect.applyEffect(".toolbar", {
	lightColor: "rgba(255,255,255,0.1)",
	gradientSize: 500
});

FluentRevealEffect.applyEffect(".toolbar > .btn", {
	clickEffect: true
});

FluentRevealEffect.applyEffect(".effect-group-container", {
	clickEffect: true,
	lightColor: "rgba(255,255,255,0.6)",
	gradientSize: 80,
	isContainer: true,
	children: {
		borderSelector: ".btn-border",
		elementSelector: ".btn",
		lightColor: "rgba(255,255,255,0.3)",
		gradientSize: 150
	}
});
