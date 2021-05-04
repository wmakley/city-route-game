import './admin.css'
import Turbolinks from 'turbolinks'
import { Application } from "stimulus"
import { definitionsFromContext } from "stimulus/webpack-helpers"
import { default as BootstrapModal } from 'bootstrap/js/dist/modal';

Turbolinks.start()

const application = Application.start()
const context = require.context("./controllers", true, /\.js$/)
application.load(definitionsFromContext(context))

const domParser = new DOMParser()

const Modal = {
	element: null,
	title: null,
	body: null,
	okayButton: null,
	okayListener: null,
	bootstrapModal: null,

	show(title, body, onOkay) {
		if (!this.element) {
			this.element = document.getElementById("defaultModal")
			this.title = this.element.querySelector(".modal-title")
			this.body = this.element.querySelector(".modal-body")
			this.okayButton = this.element.querySelector(".btn-primary")
			this.bootstrapModal = new BootstrapModal(this.element)
		}

		this.title.innerHTML = title
		this.body.innerHTML = body

		if (this.okayListener) {
			this.okayButton.removeEventListener("click", this.okayListener)
			this.okayListener = null
		}

		this.okayButton.addEventListener("click", onOkay)
		this.okayListener = onOkay

		this.bootstrapModal.show()
	},

	dispose() {
		if (!this.element) {
			return
		}

		if (this.okayListener) {
			this.okayButton.removeEventListener("click", this.okayListener)
			this.okayListener = null
		}
		this.element = null
		this.title = null
		this.body = null
		this.okayButton = null
		this.bootstrapModal.dispose()
		this.bootstrapModal = null
	}
}

let boardEditor = null;

document.addEventListener("turbolinks:before-cache", function() {
	Modal.dispose()
	if (boardEditor) {
		boardEditor.unmount()
		boardEditor = null
	}
})

document.addEventListener("turbolinks:load", function() {

	document.querySelectorAll("a[data-confirm], a[data-method]").forEach((link) => {
		const confirm = link.getAttribute("data-confirm")
		const method = link.getAttribute("data-method")

		link.addEventListener("click", (event) => {
			event.preventDefault()

			if (confirm) {
				// show bootstrap modal first
				Modal.show("Confirm", confirm, () => specialLinkNavigation(link))
				return
			} else {
				// navigate immediately
				specialLinkNavigation(link)
			}
		})
	})

	document.querySelectorAll("form[data-remote]").forEach((form) => {
		if (!form.id) {
			throw new Error("remote form must have an id for replacement to work")
		}

		form.addEventListener('submit', async function (event) {
			event.preventDefault()

			form.querySelectorAll(`input[type="submit"], button[type="submit"]`).forEach((submitBtn) => {
				submitBtn.disabled = true
			})

			const formData = new FormData(form);

			const response = await fetch(form.action, {
				method: form.method,
				body: formData,
				headers: {
					"Accept": "text/html, text/javascript",
					"X-Requested-With": "XMLHttpRequest",
				},
				mode: "same-origin",
			})
			const body = await response.text()

			if (response.ok && response.headers.get("Content-Type").startsWith("text/javascript")) {
				// Execute turbolinks visit javascript
				eval(body)
			} else {
				// Replace form
				const page = domParser.parseFromString(body, "text/html")
				form.innerHTML = page.getElementById(form.id).innerHTML
			}
		})
	})

	const eltToFocus = document.querySelector("[data-focus]")
	if (eltToFocus) {
		eltToFocus.focus()
		if ((eltToFocus.tagName === "INPUT" && eltToFocus.type === "text") || eltToFocus.tagName === "TEXTAREA") {
			eltToFocus.setSelectionRange(eltToFocus.value.length, eltToFocus.value.length)
		}
	}

	const boardEditorDiv = document.getElementById("board-editor")
	if (boardEditorDiv) {
		const boardId = parseInt(boardEditorDiv.getAttribute("data-board-id"), 10);
		if (isNaN(boardId)) {
			throw new Error("data-board-id is not a number");
		}

		import('./board-editor')
			.then(({ default: createApp }) => {
				boardEditor = createApp(boardId)
				boardEditor.mount(boardEditorDiv)
			})
			.catch((error) => {
				console.error(error)
				window.alert('An error occurred while loading the board editor \u{1F622}:\n\n' + error)
			})
	}
})

function specialLinkNavigation(link) {
	const method = link.getAttribute("data-method")

	if (method) {
		(async () => {
			const response = await fetch(link.href, {
				method: method.toUpperCase(),
				headers: {
					"Accept": "text/javascript",
					"X-Requested-With": "XMLHttpRequest"
				}
			})

			const body = await response.text()

			if (!response.ok) {
				window.alert(body)
				return
			}

			// Assume javascript body
			eval(body)
		})()
	} else {
		window.location.href = link.href
	}
}
