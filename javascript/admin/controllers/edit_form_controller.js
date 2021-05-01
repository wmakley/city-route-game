import { Controller } from "stimulus"
import Turbolinks from "turbolinks"

export default class extends Controller {
	static targets = [
		"form",
		"renameButton",
		"submitButton",
		"cancelButton",
		"nameInput",
	]

	connect() {
		// Hack to allow server to call a method in this controller
		window.updateFormSucceeded = () => {
			this.hideForm()
			Turbolinks.visit(window.location.href, {
				mode: "replace"
			})
		}
	}

	disconnect() {
		window.updateFormSucceeded = null
		this.hideForm()
	}

	get formIsVisible() {
		return this.formTarget.style.display != 'none'
	}

	hideForm() {
		if (this.formIsVisible) {
			this.toggle()
		}
	}

	toggle() {
		if (this.formIsVisible) {
			this.renameButtonTarget.style.display = ''
			this.formTarget.style.display = 'none'
		} else {
			this.renameButtonTarget.style.display = 'none'
			this.formTarget.style.display = 'block'
			this.submitButtonTarget.disabled = false
			this.nameInputTarget.focus()
			this.nameInputTarget.setSelectionRange(0, this.nameInputTarget.value.length)
		}
	}

	cancelButtonClick(event) {
		event.preventDefault()
		this.hideForm()
	}
}
