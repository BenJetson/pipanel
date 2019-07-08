/**
 *  PiPanel
 *
 *  Copyright 2019 Ben Godfrey
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not 
 *  use this file except in compliance with the License. You may obtain a copy 
 *  of the License at:
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software 
 *  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT 
 *  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the 
 *  License for the specific language governing permissions and limitations 
 *  under the License.
 */

import groovy.json.JsonOutput
import java.net.URI

metadata {
	definition (
		name: "PiPanel", 
		namespace: "BenJetson", 
		author: "Ben Godfrey", 
		cstHandler: true) {
			capability "Alarm"
			capability "Chime"
			capability "Notification"
			capability "Speech Synthesis"
			capability "Tone"
			command    "chime"
			command    "reboot"
	}
	preferences {
		input(
			"ip",
			"string",
			title: "Local IP",
			description: "PiPanel Server Address",
			required: true
		)
		input(
			"port",
			"string",
			title: "Port",
			description: "PiPanel Server Port",
			required: true
		)
	}

	simulator {
		// TODO: define status and reply messages here
	}

	tiles {
		// TODO: define your main and details tiles here
		standardTile("tone", "device.tone", decoration: "flat") {
            state "default", 
			label: "Tone", 
			action: "tone.beep", 
			icon: "st.alarm.beep.beep"
        }
		standardTile("chime", "device.tone", decoration: "flat") {
            state "default", 
			label: "Chime", 
			action: "tone.chime", 
			icon: "st.Food & Dining.dining8"
        }
	}
}

/**
 * ApiDo creates a HubAction instance that will hit the PiPanel API on the 
 * local device.
 * 
 * @param data     The object that will be marshalled to JSON and sent 
 *                 to the API.
 */
private makeApiHubAction(String actionType, Map data) {	
	return new physicalgraph.device.HubAction(
		method: "POST",
		path: "/" + actionType,
		headers: [
			"HOST":  ip + ":" + port,
			"Content-Type": "application/json"
		],
		body: JsonOutput.toJson(data)
	)
}

// parse events into attributes
def parse(String description) {
	log.debug "Parsing '${description}'"
	// TODO: handle 'alarm' attribute
	// TODO: handle 'chime' attribute

}

// handle commands
def off() {
	log.debug "Executing 'off'"
	// TODO: handle 'off' command
}

def strobe() {
	log.debug "Executing 'strobe'"
	// TODO: handle 'strobe' command
}

def siren() {
	log.debug "Executing 'siren'"
	// TODO: handle 'siren' command
}

def both() {
	log.debug "Executing 'both'"
	// TODO: handle 'both' command
}

def chime() {
	log.debug "Executing 'chime'"

	return makeApiHubAction("sound", [
		"tone": "chime"
	])
}

def beep() {
	log.debug "Executing 'beep'"
	
	return makeApiHubAction("sound", [
		"tone": "beep"
	])
}

private makeAlertParameters(String text) {
	def splitText = text.split("&")
	def numParams = splitText.size()

	if (numParams < 1) {
		return [:]
	}

	def parameters = [:]

	parameters.put("message", splitText[0])

	for (int i=1;i<numParams; i++) {
		def pair = splitText[i].split("=")

		if (pair.size() != 2) continue
		
		def key = pair[0].toLowerCase()
		def value = new URI(pair[1]).getPath()

		// Try to cast to the most relevant type possible.
		if (value.isInteger()) {
			parameters.put(key, value as Integer)
		} else if (value.isNumber()) {
			parameters.put(key, value as Double)
		} else if (value.toLowerCase() == "true") {
			parameters.put(key, true)
		} else if (value.toLowerCase() == "false") {
			parameters.put(key, false)
		} else {
			parameters.put(key, value)
		}
	} 

	return parameters
}

def speak(text) {
	log.debug "Executing 'speak'"
	
	return makeApiHubAction("alert", makeAlertParameters(text))
}

def deviceNotification(text) {
	log.debug "Executing 'deviceNotification'"
	
	return makeApiHubAction("alert", makeAlertParameters(text))
}

// HubAction Helper Methods
// These below are adapted from the methods provided by SmartThings. See also:
// ST Classic Docs > Cloud and LAN Connected Devices > Building the Device Type

/**
 * getCallbackAddress gets the address of the Hub.
 */
private getCallbackAddress() {
    return device.hub.getDataValue("localIP") + ":" + 
		device.hub.getDataValue("localSrvPortTCP")
}