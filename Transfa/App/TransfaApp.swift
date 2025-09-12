/**
 * @description
 * This file contains the main entry point for the Transfa SwiftUI application.
 *
 * It defines the `TransfaApp` struct which conforms to the `App` protocol.
 * The ` @main` attribute designates this as the starting point of the application execution.
 * The body of the app scene contains a `WindowGroup`, which is the container for the app's UI.
 *
 * Initially, it presents a simple `ContentView` as the root view. This will be expanded
 * later to include environment objects, state management, and routing logic to handle
 * authentication status and user onboarding.
 *
 * @dependencies
 * - SwiftUI: The core framework for building the application's user interface.
 */
import SwiftUI

@main
struct TransfaApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView()
        }
    }
}

/**
 * @description
 * A temporary placeholder view that will be replaced by the main app router.
 *
 * This view is used to ensure the project compiles and runs from the start. It will be
 * removed in a future step when the onboarding and home screen features are implemented.
 */
struct ContentView: View {
    var body: some View {
        VStack {
            Image(systemName: "dollarsign.circle.fill")
                .imageScale(.large)
                .foregroundColor(.accentColor)
            Text("Welcome to Transfa!")
        }
        .padding()
    }
}